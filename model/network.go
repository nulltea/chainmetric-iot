package model

type Network struct {
	IPAddress  string `json:"ip"`
	MACAddress string `json:"mac"`
	Hostname   string `json:"hostname"`
}
