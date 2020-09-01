package config

type Info struct {
	Url             string `json:"url"`
	PostEdge        string `json:"postEdge"`
	Socket          string `json:"socket"`
	PutStatus       string `json:"putStatus"`
	GetInfo         string `json:"getInfo"`
	Key             string `json:"key"`
	PostDevice      string `json:"postDevice"`
	WebsocketServer string `json:"websocketServer"`
	//demo	GetCommand   string `json:"getCommand"`
}
