OUTPUT=bin/sensor
REMOTE_IP=192.168.31.180
REMOTE_DIR=/home/pi/sensorsys
CRYPTO_DIR=../network/crypto-config/peerOrganizations/supplier.iotchain.network/users/User1@supplier.iotchain.network/msp

build:
	CGO_ENABLED=1 go mod vendor && go build -v  -o $(OUTPUT) .

build-remote:
	CC=arm-linux-gnueabihf-gcc CXX=arm-linux-gnueabihf-g++ \
        CGO_ENABLED=1 GOOS=linux GOARCH=arm64 GOARM=6 \
        go build -v  -o $(OUTPUT) .

sync:
	rsync -r --exclude .env \
	 --delete --filter 'P vendor' --filter 'P bin' \
	. pi@$(REMOTE_IP):$(REMOTE_DIR)

crypto-sync:
	scp $(CRYPTO_DIR)/signcerts/User1@supplier.iotchain.network-cert.pem pi@$(REMOTE_IP):identity.pem
	scp $(CRYPTO_DIR)/keystore/priv_sk pi@$(REMOTE_IP):identity.key

run:
	sudo ./$(OUTPUT)

kill:
	ps aux |  awk '{print $2"\t"$11}' | grep -E '^\d+\t'"\."'$' | awk '{print $1}' | xargs kill -SIGTERM

i2c:
	sudo i2cdetect -l
	sudo i2cdetect -y 1
	sudo i2cdetect -y 2
	sudo i2cdetect -y 3
	sudo i2cdetect -y 4
