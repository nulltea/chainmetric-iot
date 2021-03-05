module sensorsys

go 1.15

require (
	github.com/cgxeiji/max3010x v0.0.0-20200914015011-b05e3d2950ea
	github.com/d2r2/go-bsbmp v0.0.0-20190515110334-3b4b3aea8375
	github.com/d2r2/go-dht v0.0.0-20200119175940-4ba96621a218
	github.com/d2r2/go-i2c v0.0.0-20191123181816-73a8a799d6bc
	github.com/d2r2/go-logger v0.0.0-20181221090742-9998a510495e
	github.com/d2r2/go-shell v0.0.0-20191113051817-7664ea33645f // indirect
	github.com/hyperledger/fabric-sdk-go v1.0.0
	github.com/op/go-logging v0.0.0-20160315200505-970db520ece7
	gopkg.in/yaml.v2 v2.4.0
)

replace github.com/cgxeiji/max3010x v0.0.0-20200914015011-b05e3d2950ea => github.com/timoth-y/max3010x v0.0.0-20210227142301-7a3d7c5be5c7
