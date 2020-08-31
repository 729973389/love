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
		log.Fatal(errors.Wrap(token.Error(),"connect broker"))
	}
	//test////////////////////////////////////////////////////////////////////////////
	c.Client.Subscribe("easyfetch/device/properties/"+"lxd", 0, messageHandler)
	//////////////////////////////////////////////////////////////////////////////////
	go c.Receive()
}

func (c *Client) Receive() {
	for {
		select {
		case message := <-c.Hub.Down:
			log.Println(string(message))
			deviceId := FindDeviceId(message)
			if deviceId == "" {
				continue
			}
			if _, ok := c.Map[deviceId]; ok {
				log.Warning(errors.Wrap(fmt.Errorf(deviceId), "  already exit"))
				continue
			}
			c.Map[deviceId] = true
			c.Client.Subscribe("easyfetch/device/properties/"+deviceId, 0, messageHandler)
		}
	}
}

var messageHandler mqtt.MessageHandler = func(client mqtt.Client, message mqtt.Message) {
	/////////c.Hub.Up <- message.Payload()
	log.Println(string(message.Payload()))
}

func FindDeviceId(b []byte) string {
	message := string(b)
	ft1 := strings.Split(message, ",")
	for _, v := range ft1 {
		if strings.Contains(v, "\"deviceId\":") {
			tokens := strings.Split(v, "\"")
			for i2, v2 := range tokens {
				if v2 == ":" {
					return tokens[i2+1]
				}
			}
		}
	}
	log.Error("MQTT can‘t find deviceId")
	return ""
}

func (c *Client) PublishCommand() {
	for {
		select {
		case command := <-c.Hub.Command:
			deviceId := FindDeviceId(command)
			if token := c.Client.Publish("easyfetch/device/command/"+deviceId, 1, false, command); token.Wait() && token.Error() != nil {
				log.Warning(token.Error())
			}
		}
	}
}
