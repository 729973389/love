package mqtter

import (
	mqtt"github.com/eclipse/paho.mqtt.golang"
)

type Client struct {
	Client *mqtt.Client
}
