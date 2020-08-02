package main

import (
	"crypto/hmac"
	"crypto/sha256"
	"crypto/tls"
	"encoding/hex"
	"encoding/json"
	"fmt"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"log"
	"strings"
	"sync"
	"time"
)

const (
	deviceId     = "5f0bae9b52045702fceed061_myimg"
	secret       = "12345678"
	Host         = "tls://192.168.32.189:7883"
	flag         = "_"
	publishTopic = "$oc/devices/" + deviceId +
		"/sys/properties/report"
	serviceId      = "Image"
	requestId      = "CMD"
	subscribeTopic = "$oc/devices/" + deviceId +
		"/sys/commands/#"
)


type Message struct {
	Services []Services `json:"services"`
}
type Services struct {
	Service_id string     `json:"service_id"`
	Properties Properties `json:"properties"`
	EventTime  string     `json:"event_time"`
}
type Properties struct {
	Frame string `json:"frame"`
}

type Commands struct {
	ObjectDeviceID string `json:"object_device_id"`
	CommandName    string `json:"command_name"`
	ServiceID      string `json:"service_id"`
	Paras          *Paras `json:"paras"`
}
type Paras struct {
	Capture int `json:"capture"`
	Control int `json:"control"`
}

var messageHandler mqtt.MessageHandler = func(client mqtt.Client, message mqtt.Message) {
	commands := new(Commands)
	log.Println("Im handling message")
	fmt.Println(string(message.Payload()))
	err := json.Unmarshal(message.Payload(), commands)
	if err != nil {
		panic(err)
	}
	fmt.Println(commands.ObjectDeviceID)
	fmt.Println(commands.ServiceID)
	fmt.Println(commands.CommandName)
	fmt.Println(commands.Paras.Capture)
	fmt.Println(commands.Paras.Control)

}

func subscrebe(c *mqtt.Client) {
	token := (*c).Subscribe(subscribeTopic, 0, messageHandler)
	if token.Wait() && token.Error() != nil {
		panic(token.Error())
	}
	log.Println("Subscribe : " + subscribeTopic + ":OK")
}

func main() {

	var wg sync.WaitGroup
	ch := make(chan int)
	wg.Add(2)
	sha, yyyymmddhh := getSha256()
	clientId := deviceId + flag + "0" + flag + "1" + flag + yyyymmddhh
	fmt.Println(clientId)
	ops := mqtt.NewClientOptions().
		SetClientID(clientId).
		SetUsername(deviceId).
		AddBroker(Host).
		SetPassword(sha).
		SetTLSConfig(&tls.Config{InsecureSkipVerify: true}).
		SetKeepAlive(60 * time.Second).
		SetCleanSession(true)

	ops.ProtocolVersion = 4
	client := mqtt.NewClient(ops)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())

	} else {
		fmt.Println("CONNECTED!")
	}
	go func() {
		defer wg.Done()
		subscrebe(&client)
	}()
	go func() {
		defer wg.Done()
		for  {
			pop(&client)
			switch <-ch {
			case 1:
				log.Println("******STOPING SENDING******")
				return
			}
		}

	}()

	wg.Wait()
	fmt.Println("Im going to close it!")
	time.Sleep(5 * time.Second)

}

func getSha256() (s string, t string) {
	var yyyymmddhh string
	format := "2006-01-02 15:04:05"
	utc := time.Now().UTC().Format(format)
	utc = strings.Replace(utc, "-", "", -1)
	utc = strings.Replace(utc, " ", "", -1)
	utc = strings.Replace(utc, ":", "", -1)
	for i, v := range utc {
		if i < 10 {
			yyyymmddhh += string(v)
		}
	}
	fmt.Println("TIME FLAG: " + yyyymmddhh)
	sha := hmac.New(sha256.New, []byte(yyyymmddhh))
	sha.Write([]byte(secret))
	fmt.Println(hex.EncodeToString(sha.Sum(nil)))
	return hex.EncodeToString(sha.Sum(nil)), yyyymmddhh

}

func mashalMessage() interface{} {
	p := &Properties{
		Frame: "/mnt/hgfs/wuff/gcv/" + getEventTime() + ".jpg",
	}

	s := &Services{
		Service_id: serviceId,
		Properties: *p,
		EventTime:  getEventTime(),
	}
	m := Message{
		Services: []Services{*s},
	}
	marshalMessage, err := json.Marshal(m)
	if err != nil {
		fmt.Println(err)

	}
	fmt.Println(string(marshalMessage))
	return string(marshalMessage)
}

func getEventTime() string {
	var evenTime string
	format := "2006-01-02 15:04:05"
	utc := time.Now().UTC().Format(format)
	utc = strings.Replace(utc, "-", "", -1)
	utc = strings.Replace(utc, " ", "", -1)
	utc = strings.Replace(utc, ":", "", -1)
	for i, v := range utc {
		evenTime += string(v)
		if i == 7 {
			evenTime += "T"
		}
		if i == len(utc)-1 {
			evenTime += "Z"
		}
	}
	fmt.Println(evenTime)
	return evenTime
}

func pop(client *mqtt.Client) {
	if token := (*client).Publish(publishTopic, 0, false, mashalMessage());
		token.Wait() && token.Error() != nil {
		panic(token.Error())
	} else {
		fmt.Println("PUBLISHED")
	}
}
