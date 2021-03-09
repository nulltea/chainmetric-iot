package model

type DeviceSignature struct {
	Network
	Supports []string `json:"supports"`
}
