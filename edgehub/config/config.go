package config

type Info struct {
	EdgeInfoServer      string `json:"edgeInfoServer"`
	PostEdge            string `json:"postEdge"`
	Socket              string `json:"socket"`
	PutStatus           string `json:"putStatus"`
	GetInfo             string `json:"getInfo"`
	Key                 string `json:"key"`
	PostDevice          string `json:"postDevice"`
	DeviceInfoServer    string `json:"deviceInfoServer"`
	DeviceControlServer string `json:"deviceControlServer"`
}
