package config

type Url struct {
	Url       string `json:"url"`
	SendData  string `json:"sendData"`
	Socket    string `json:"socket"`
	PutStatus string `json:"putStatus"`
	GetInfo string `json:"getInfo"`
}
