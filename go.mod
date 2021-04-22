module github.com/timoth-y/iot-blockchain-sensorsys

go 1.15

require (
	github.com/bskari/go-lsm303 v0.0.0-20200927082938-3432d22cb4f1
	github.com/cgxeiji/max3010x v0.0.0-20200914015011-b05e3d2950ea
	github.com/d2r2/go-dht v0.0.0-20200119175940-4ba96621a218
	github.com/d2r2/go-logger v0.0.0-20181221090742-9998a510495e
	github.com/d2r2/go-shell v0.0.0-20191113051817-7664ea33645f // indirect
	github.com/hyperledger/fabric-sdk-go v1.0.0
	github.com/op/go-logging v0.0.0-20160315200505-970db520ece7
	github.com/pkg/errors v0.9.1
	github.com/skip2/go-qrcode v0.0.0-20200617195104-da1b6568686e
	github.com/timoth-y/iot-blockchain-contracts v0.0.0-20210316013042-82c012a55cea
	gopkg.in/yaml.v2 v2.4.0
	periph.io/x/periph v3.6.7+incompatible
)

replace github.com/cgxeiji/max3010x v0.0.0-20200914015011-b05e3d2950ea => github.com/timoth-y/max3010x v0.0.0-20210310203014-cf62a2a2aea3

replace github.com/bskari/go-lsm303 v0.0.0-20200927082938-3432d22cb4f1 => github.com/timoth-y/go-lsm303 v0.0.0-20210418113746-a26cc847773d
