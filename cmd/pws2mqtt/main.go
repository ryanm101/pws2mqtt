package main

import (
	"fmt"
	"github.com/eclipse/paho.mqtt.golang"
	"github.com/mitchellh/mapstructure"
	"net/http"
	"os"
	"time"
)

var messagePubHandler mqtt.MessageHandler = func(client mqtt.Client, msg mqtt.Message) {
	fmt.Printf("Message %s received on topic %s\n", msg.Payload(), msg.Topic())
}

var connectHandler mqtt.OnConnectHandler = func(client mqtt.Client) {
	fmt.Println("Connected")
}

var connectionLostHandler mqtt.ConnectionLostHandler = func(client mqtt.Client, err error) {
	fmt.Printf("Connection Lost: %s\n", err.Error())
}

func setupMQTT(broker string, user string, pass string) mqtt.Client {
	options := mqtt.NewClientOptions()
	options.Password = pass
	options.Username = user
	options.AddBroker(broker)
	options.SetClientID("pws2mqtt")
	options.SetDefaultPublishHandler(messagePubHandler)
	options.OnConnect = connectHandler
	options.OnConnectionLost = connectionLostHandler

	return mqtt.NewClient(options)
}

func ecowittHandler(pwsdata Pwsdata, mqttclient mqtt.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		headerContentType := r.Header.Get("Content-Type")
		if headerContentType != "application/x-www-form-urlencoded" {
			w.WriteHeader(http.StatusUnsupportedMediaType)
			return
		}
		r.ParseForm()
		form := make(map[string]interface{})
		for key, value := range r.Form {
			//fmt.Printf("Key:%s, Value:%s\n", key, value[0])
			form[key] = value[0]
		}
		mapstructure.Decode(form, &pwsdata)
		sendMQTTUpdate(mqttclient, pwsdata.ToJSON())
		return
	}
}

func wundergroundHandler(pwsdata Pwsdata, mqttclient mqtt.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		values := r.URL.Query()
		form := make(map[string]interface{})
		for key, value := range values {
			//fmt.Printf("Key:%s, Value:%s\n", key, value[0])
			form[key] = value[0]
		}
		mapstructure.Decode(form, &pwsdata)
		sendMQTTUpdate(mqttclient, pwsdata.ToJSON())
		return
	}
}

func sendMQTTUpdate(mqttclient mqtt.Client, data string) {
	token := mqttclient.Connect()
	if token.Wait() && token.Error() != nil {
		panic(token.Error())
	}

	LWT := "pws/tele/LWT"
	State := "pws/tele/STATE"
	Sensor := "pws/tele/SENSOR"
	LWTtext := fmt.Sprintf("%s", "Online")
	Statetext := fmt.Sprintf("%s", "{}")
	Sensortext := fmt.Sprintf("%s", data)
	token = mqttclient.Publish(LWT, 0, false, LWTtext)
	token.Wait()
	token = mqttclient.Publish(State, 0, false, Statetext)
	token.Wait()
	token = mqttclient.Publish(Sensor, 0, false, Sensortext)
	token.Wait()
	time.Sleep(time.Second)
	mqttclient.Disconnect(100)
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

func main() {
	listenIP := getEnv("IPADDR", "")
	listenPort := getEnv("LISTENPORT", "8080")
	mqttServer := getEnv("MQTTSERVER", "test.mosquitto.org")
	mqttPort := getEnv("MQTTPORT", "1883")
	mqttUser := getEnv("MQTTUSER", "")
	mqttPass := getEnv("MQTTPASS", "")

	var edata Ecowittdata
	var wdata Wundergrounddata

	mqttclient := setupMQTT(fmt.Sprintf("tcp://%s:%s", mqttServer, mqttPort),
		mqttUser, mqttPass)

	http.Handle("/ecowitt", ecowittHandler(&edata, mqttclient))
	http.Handle("/wunderground", wundergroundHandler(&wdata, mqttclient))

	http.ListenAndServe(fmt.Sprintf("%s:%s", listenIP, listenPort), nil)
}
