package main

import (
	"fmt"
	"github.com/eclipse/paho.mqtt.golang"
	"github.com/mitchellh/mapstructure"
	log "github.com/sirupsen/logrus"
	"net/http"
	"os"
	"time"
)

var messagePubHandler mqtt.MessageHandler = func(client mqtt.Client, msg mqtt.Message) {
	log.Debugf("Message %s received on topic %s\n", msg.Payload(), msg.Topic())
}

var connectHandler mqtt.OnConnectHandler = func(client mqtt.Client) {
	log.Infof("MQTT Connected")
}

var connectionLostHandler mqtt.ConnectionLostHandler = func(client mqtt.Client, err error) {
	log.Errorf("MQTT Connection Lost: %s\n", err.Error())
}

func setupMQTT(broker string, user string, pass string) mqtt.Client {
	log.Debug("setupMQTT")
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
		log.Info("wundergroundHandler")
		headerContentType := r.Header.Get("Content-Type")
		if headerContentType != "application/x-www-form-urlencoded" {
			w.WriteHeader(http.StatusUnsupportedMediaType)
			return
		}
		r.ParseForm()
		form := make(map[string]interface{})
		for key, value := range r.Form {
			log.Debugf("Key:%s, Value:%s\n", key, value[0])
			form[key] = value[0]
		}
		mapstructure.Decode(form, &pwsdata)
		log.Debug(pwsdata.ToJSON())
		sendMQTTUpdate(mqttclient, pwsdata.ToJSON())
		return
	}
}

func wundergroundHandler(pwsdata Pwsdata, mqttclient mqtt.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Info("wundergroundHandler")
		values := r.URL.Query()
		form := make(map[string]interface{})
		for key, value := range values {
			log.Debugf("Key:%s, Value:%s\n", key, value[0])
			form[key] = value[0]
		}
		mapstructure.Decode(form, &pwsdata)
		log.Debug(pwsdata.ToJSON())
		sendMQTTUpdate(mqttclient, pwsdata.ToJSON())
		return
	}
}

func sendMQTTUpdate(mqttclient mqtt.Client, data string) {
	log.Debug("sendMQTTUpdate")
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

func init() {
	// Log as JSON instead of the default ASCII formatter.
	log.SetFormatter(&log.JSONFormatter{})

	// Output to stdout instead of the default stderr
	// Can be any io.Writer, see below for File example
	log.SetOutput(os.Stdout)

	// Only log the warning severity or above.
	log.SetLevel(log.WarnLevel)
}

func main() {
	log.SetFormatter(&log.JSONFormatter{})
	log.SetOutput(os.Stdout)
	log.SetLevel(log.InfoLevel)

	if getEnv("DEBUG", "") == "true" {
		log.SetLevel(log.DebugLevel)
	}

	listenIP := getEnv("IPADDR", "")
	listenPort := getEnv("LISTENPORT", "8080")
	mqttServer := getEnv("MQTTSERVER", "")
	mqttPort := getEnv("MQTTPORT", "1883")
	mqttUser := getEnv("MQTTUSER", "")
	mqttPass := getEnv("MQTTPASS", "")

	var edata Ecowittdata
	var wdata Wundergrounddata

	mqttclient := setupMQTT(fmt.Sprintf("tcp://%s:%s", mqttServer, mqttPort),
		mqttUser, mqttPass)

	http.Handle("/ecowitt", ecowittHandler(&edata, mqttclient))
	http.Handle("/wunderground", wundergroundHandler(&wdata, mqttclient))

	log.Infof("Listening On: %s:%s", listenIP, listenPort)
	http.ListenAndServe(fmt.Sprintf("%s:%s", listenIP, listenPort), nil)
}
