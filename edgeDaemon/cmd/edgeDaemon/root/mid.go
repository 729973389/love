package root

import (
	"crypto/tls"
	"fmt"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"strings"
)

type Client struct {
	Hub    *Hub
	Client mqtt.Client
	Send   chan []byte
	Map    map[string]bool
}

var c *Client

//func NewClient()*Client{
//	ops := mqtt.NewClientOptions().SetClientID("edgeDaemon").
//		SetUsername("edgeDaemon").
//		SetPassword("12345678").
//		SetTLSConfig(&tls.Config{InsecureSkipVerify: true}).
//		SetAutoReconnect(true).
//		SetProtocolVersion(4).
//		SetCleanSession(true)
//	client := mqtt.NewClient(ops)
//	return &Client{Client: &client}
//}

func Run(hub *Hub) {
	ops := mqtt.NewClientOptions().SetClientID("edgeDaemon").AddBroker("tcp://localhost:7883").
		SetUsername("edgeDaemon").
		SetPassword("12345678").
		SetTLSConfig(&tls.Config{InsecureSkipVerify: true}).
		SetAutoReconnect(true).
		SetProtocolVersion(4).
		SetCleanSession(true)
	client := mqtt.NewClient(ops)
	c = &Client{Hub: hub, Client: client, Send: make(chan []byte, 1024), Map: make(map[string]bool)}
	if token := c.Client.Connect(); token.Wait() && token.Error() != nil {
		log.Fatal(errors.Wrap(token.Error(), "connect broker"))
	}
	//test////////////////////////////////////////////////////////////////////////////
	c.Client.Subscribe("easyfetch/device/properties/"+"lxd", 0, messageHandler)
	//////////////////////////////////////////////////////////////////////////////////
	go c.Receive()
	//go c.PublishCommand()
}

func (c *Client) Receive() {
	for {
		select {
		case b := <-c.Hub.Down:
			s := string(b)
			deviceId := FindKeyString(s, "deviceId")
			if deviceId == "" {
				continue
			}
			switch t := FindKeyString(s, "type"); t {
			case "bind":
				log.Info("bind")
				if _, ok := c.Map[deviceId]; ok {
					log.Warning(errors.Wrap(fmt.Errorf(deviceId), "already exit"))
					continue
				}
				c.Map[deviceId] = true
				c.Client.Subscribe("easyfetch/device/properties/"+deviceId, 0, messageHandler)
			case "unbind":
				log.Info("unbind")
				if _, ok := c.Map[deviceId]; !ok {
					log.Warning("No such deviceId")
					continue
				}
				c.Client.Unsubscribe("easyfetch/devices/properties" + deviceId)
			}



		}
	}
}

var messageHandler mqtt.MessageHandler = func(client mqtt.Client, message mqtt.Message) {
	c.Hub.Up <- message.Payload()
	log.Println(string(message.Payload()))
}

func FindKeyString(s, k string) string {
	ft1 := strings.Split(s, ",")
	for _, v := range ft1 {
		if strings.Contains(v, "\""+k+"\":") {
			tokens := strings.Split(v, "\"")
			for i2, v2 := range tokens {
				if strings.Contains(v2, ":") {
					if tokens[i2-1] == k {
						return tokens[i2+1]
					}
				}
			}
		}
	}
	log.Error("MQTT can‘t find ", k)
	return ""
}

func (c *Client) PublishCommand() {
	for {
		select {
		case command := <-c.Hub.Command:
			/////////////command
			deviceId := FindKeyString(string(command), "co")
			if token := c.Client.Publish("easyfetch/device/command"+deviceId, 1, false, command); token.Wait() && token.Error() != nil {
				log.Warning(token.Error())
				continue
			}
			c.Client.Subscribe("easyfetch/device/"+deviceId+"/command/response", 0, messageHandler)
		}
	}
}
