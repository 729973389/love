package root

import (
	"context"
	"crypto/tls"
	"fmt"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"strings"
	"sync"
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

func RunMQTT(ctx context.Context, hub *Hub) {
	defer log.Warning("EXIT RUNMQTT")
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
		log.Error(errors.Wrap(token.Error(), "connect broker"))
		return
	}
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		c.Receive(ctx)
	}()
	//go c.PublishCommand()
	select {
	case <-ctx.Done():
		break
	}
	wg.Wait()
}

func (c *Client) Receive(ctx context.Context) {
	defer func() {
		log.Warning("EXIT RECEIVE")
		c.Client.Disconnect(200)
	}()
	for {
		select {
		case deviceId := <-c.Hub.DeviceMap:
			log.Info("register deviceId")
			for _, v := range deviceId {
				c.Map[v] = true
				c.Client.Subscribe("easyfetch/device/properties/"+v, 0, messageHandler)
			}

		case deviceGister := <-c.Hub.Down:
			log.Info("receive: ", deviceGister.String())
			deviceId := deviceGister.DeviceId
			if deviceId == "" {
				continue
			}
			switch t := deviceGister.Type; t {
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
				c.Client.Unsubscribe("easyfetch/device/properties/" + deviceId)
				delete(c.Map, deviceId)
			}
		case <-ctx.Done():
			return
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
			deviceId := FindKeyString(string(command), "deviceId")
			if token := c.Client.Publish("easyfetch/device/command"+deviceId, 1, false, command); token.Wait() && token.Error() != nil {
				log.Warning(token.Error())
				continue
			}
			c.Client.Subscribe("easyfetch/device/"+deviceId+"/command/response", 0, messageHandler)
		}
	}
}
