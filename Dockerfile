# syntax=docker/dockerfile:1
FROM golang:1.18 as builder
WORKDIR /src
COPY . ./
RUN make build-native

FROM alpine:latest
ENV LISTENIP=""
ENV LISTENPORT="8080"
ENV MQTTSERVER="test.mqtt.org"
ENV MQTTPORT="1883"
ENV MQTTUSER=""
ENV MQTTPASS=""

EXPOSE 8080
WORKDIR /app

RUN apk --no-cache add ca-certificates && \
    addgroup -g 1000 -S pws2mqtt && \
    adduser -u 1000 -S pws2mqtt -G pws2mqtt

COPY --from=builder /src/dist/pws2mqtt /app/pws2mqtt

RUN chown -R pws2mqtt:pws2mqtt /app && \
    chmod -R 754 /app

USER pws2mqtt

ENTRYPOINT ["./pws2mqtt"]
