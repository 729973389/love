package config

type Server struct {
	Url          string `json:"url,omitempty"`
	Token        string `json:"token"`
	SerialNumber string `json:"serialNumber"`
}
