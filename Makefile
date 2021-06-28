include .env
export

OUTPUT=bin/sensorsys
REMOTE_DIR=/home/pi/sensorsys

build:
	GOOS=linux GOARCH=arm GOARM=6 \
		go build -v  -o $(OUTPUT) .

sync:
	rsync -r --exclude .env \
	 --delete --filter 'P vendor' --filter 'P bin' --filter 'P keystore' \
	. pi@$(REMOTE_IP):$(REMOTE_DIR)

setup-device: build
	$(eval orgHostname := $(ORG).org.$(DOMAIN))
	$(eval userIdentity := $(USER_ID)@$(orgHostname))
	$(eval mspPath := $(CRYPTO_DIR)/peerOrganizations/$(orgHostname)/users/$(userIdentity)/msp)

	../network/fabnctl gen connection -f ../network/network-config.yaml -n edge-device \
		-c supply-channel -o chipa-inu -x=device-userID=edge-device ../network

	scp $(mspPath)/signcerts/$(userIdentity)-cert.pem pi@$(REMOTE_IP):identity.pem
	scp $(mspPath)/keystore/priv_sk pi@$(REMOTE_IP):identity.key

	$(MAKE) sync

update-device: build sync

run:
	sudo "./$(OUTPUT)"

kill:
	ps aux | awk '{print $$2"\t"$$11}' | grep -E ./$(OUTPUT) | awk '{print $$1}' | sudo xargs kill -SIGTERM

i2c:
	sudo i2cdetect -l
	sudo i2cdetect -y 1
	sudo i2cdetect -y 3
	sudo i2cdetect -y 4
	sudo i2cdetect -y 5
