# PWS2MQTT
[![ryanm101](https://circleci.com/gh/ryanm101/pws2mqtt.svg?style=svg)]()

Captures output from ecowitt personal weather station and replays the data into MQTT. 

Supports both ecowitt format or wunderground formats

## Usage

All settings are configured through Environment Variables

```bash
LISTENIP=""
LISTENPORT="8080"
MQTTSERVER=""
MQTTPORT="1883"
MQTTUSER=""
MQTTPASS=""
```

* ``LISTENIP`` -> IP Address to listen on if multiple IPs are configured, default - Listen on all
* ``LISTENPORT`` -> Port for application to listen on - Default 8080
* ``MQTTSERVER`` -> MQTT Server to send data to - Default ""
* ``MQTTPORT`` -> MQTT Port to connect to - Default 1883
* ``MQTTUSER`` -> Username to use for MQTT auth - Default ""
* ``MQTTPASS`` -> Password to use for MQTT auth - Default ""

### Docker
* ``docker run --name pws2mqtt -d -e MQTTSERVER=test.mqtt.server -p 8080:8080 pws2mqtt``

### Docker-Compose
* Edit ``pws2mqtt.env``
* ``docker-compose up``

### Native

* ``MQTTSERVER=test.mqtt.server ./pws2mqtt``

## Build Instructions

### Build native
* Install GoLang
* ``make build-native``

### Build all archs
* ``make``
### Build Docker
* ``make docker``

## TODO
* [ ] Convert Units from raw data to localised (UK, metric, imperial)
* [ ] Send STATE data to MQTT
* [ ] Helm Chart
* [ ] HASSIO Add-on
* [ ] Build binaries assets with CICD
* [ ] Upload to dockerhub
* [ ] add DEBUG output
* [ ] Tests
* [ ] Allow Topic change