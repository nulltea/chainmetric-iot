package model

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/timoth-y/chainmetric-core/models"
)

type DeviceSpecs struct {
	Network
	Supports []models.Metric `json:"supports"`
}

func (ds DeviceSpecs) Encode() string {
	var metrics []string

	for i := range ds.Supports {
		metrics = append(metrics, string(ds.Supports[i]))
	}

	return fmt.Sprintf("${%s;%s;%s}", ds.Hostname, ds.IPAddress, strings.Join(metrics, ","))
}

func (ds DeviceSpecs) EncodeJson() string {
	b, _ := json.Marshal(ds)
	return string(b)
}
