package model

import (
	"encoding/json"
	"fmt"
	"strings"
)

type DeviceSpecs struct {
	Network
	Supports []string `json:"supports"`
}

func (ds *DeviceSpecs) Encode() string {
	return fmt.Sprintf("${%s;%s;%s}", ds.Hostname, ds.IPAddress, strings.Join(ds.Supports, ","))
}

func (ds *DeviceSpecs) EncodeJson() string {
	b, _ := json.Marshal(*ds)
	return string(b)
}
