include .env
export

OUTPUT=bin/sensorsys
REMOTE_DIR=/home/pi/sensorsys

build:
	GOOS=linux GOARCH=arm64 \
		go build -v  -o $(OUTPUT) .

sync:
	rsync -r --exclude .env \
	 --delete --filter 'P vendor' --filter 'P bin' \
	. pi@$(REMOTE_IP):$(REMOTE_DIR)

setup-device: build
	$(eval orgHostname := "$(ORG).org.$(DOMAIN)")
	$(eval userIdentity := "$(USER_ID)@$(orgHostname)")
	$(eval mspPath := "$(CRYPTO_DIR)/peerOrganizations/$(orgHostname)/users/$(userIdentity)/msp")

	$(eval TLS_PEER_ROOT_CERT := $(shell cat "$(CRYPTO_DIR)/peerOrganizations/$(orgHostname)/tlsca/tlsca.$(orgHostname)-cert.pem"))
	$(eval TLS_ORDERER_ROOT_CERT := $(shell cat "$(CRYPTO_DIR)/ordererOrganizations/$(DOMAIN)/tlsca/tlsca.$(DOMAIN)-cert.pem"))
	$(eval TLS_CA_ROOT_CERT := $(shell cat "$(CRYPTO_DIR)/peerOrganizations/$(orgHostname)/ca/ca.$(orgHostname)-cert.pem"))

	export $TLS_PEER_ROOT_CERT
	export $TLS_ORDERER_ROOT_CERT
	export $TLS_CA_ROOT_CERT

	envsubst < connection-template.yaml > connection.yaml

	scp $(mspPath)/signcerts/$(userIdentity)-cert.pem pi@$(REMOTE_IP):identity.pem
	scp $(mspPath)/keystore/priv_sk pi@$(REMOTE_IP):identity.key

	$(MAKE) sync

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
