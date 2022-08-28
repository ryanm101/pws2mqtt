.PHONEY: docker
GOFILES = $(shell find . -name '*.go')

default: build

dist:
	mkdir -p dist

build: dist/pws2mqtt

build-native: $(GOFILES)
	go build -o dist/pws2mqtt ./cmd/pws2mqtt

dist/pws2mqtt: $(GOFILES)
	GOOS=linux   GOARCH=amd64 CGO_ENABLED=0 go build -o dist/pws2mqtt_linux_x86_64       ./cmd/pws2mqtt
	GOOS=linux   GOARCH=arm64 CGO_ENABLED=0 go build -o dist/pws2mqtt_linux_arm64        ./cmd/pws2mqtt
	GOOS=linux   GOARCH=arm   CGO_ENABLED=0 go build -o dist/pws2mqtt_linux_arm          ./cmd/pws2mqtt
	GOOS=windows GOARCH=amd64 CGO_ENABLED=0 go build -o dist/pws2mqtt_windows_x86_64.exe ./cmd/pws2mqtt
	GOOS=darwin  GOARCH=amd64 CGO_ENABLED=0 go build -o dist/pws2mqtt_darwin_x86_64      ./cmd/pws2mqtt

docker:
	docker build -t pws2mqtt:latest .

clean:
	rm -f dist/*