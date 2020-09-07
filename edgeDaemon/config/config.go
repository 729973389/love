package config

type Server struct {
	Url          string `json:"url,omitempty"`
	SerialNumber string `json:"serialNumber"`
	Token        string `json:"token"`
}
