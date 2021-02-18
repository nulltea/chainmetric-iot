OUTPUT=./bin/sensor
REMOTE_IP=192.168.31.170
REMOTE_DIR=/home/pi/sensorsys

build:
	CGO_ENABLED=1 go build -v  -o $(OUTPUT)

build-remote:
	CC=arm-linux-gnueabihf-gcc CXX=arm-linux-gnueabihf-g++ \
        CGO_ENABLED=1 GOOS=linux GOARCH=arm64 GOARM=6 \
        go build -v  -o $(OUTPUT)

upload:
	rsync -r . pi@$(REMOTE_IP):$(REMOTE_DIR)

run:
	bash $(OUTPUT)
