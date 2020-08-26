package mqtter

import (
	"crypto/tls"
	mqtt"github.com/eclipse/paho.mqtt.golang"
)

type Client struct {
	Client *mqtt.Client
	Send chan []byte


}

func NewClient()*Client{
	ops := mqtt.NewClientOptions().SetClientID("edgeDaemon").
		SetUsername("edgeDaemon").
		SetPassword("12345678").
		SetTLSConfig(&tls.Config{InsecureSkipVerify: true}).
		SetAutoReconnect(true).
		SetProtocolVersion(4).
		SetCleanSession(true)
	client := mqtt.NewClient(ops)
	return &Client{Client: &client}
}

