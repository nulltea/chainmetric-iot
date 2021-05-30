package model

// Network defines structure for storing network environment info.
type Network struct {
	IPAddress  string `json:"ip"`
	MACAddress string `json:"mac"`
	Hostname   string `json:"hostname"`
}

