OUTPUT=bin/sensor
REMOTE_IP=192.168.31.180
REMOTE_DIR=/home/pi/sensorsys

build:
	CGO_ENABLED=1 go build -v  -o $(OUTPUT)

build-remote:
	CC=arm-linux-gnueabihf-gcc CXX=arm-linux-gnueabihf-g++ \
        CGO_ENABLED=1 GOOS=linux GOARCH=arm64 GOARM=6 \
        go build -v  -o $(OUTPUT)

sync:
	rsync -r --delete . pi@$(REMOTE_IP):$(REMOTE_DIR)

run:
	sudo ./$(OUTPUT)

i2c:
	sudo i2cdetect -l
	sudo i2cdetect -y 1
	sudo i2cdetect -y 2
	sudo i2cdetect -y 3
	sudo i2cdetect -y 4
