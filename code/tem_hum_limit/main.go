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
	deviceId     = "5f0e653bc7132802fb5ed0c9_tem_humi_v1"
	secret       = "12345678"
	Host         = "tls://192.168.32.189:7883"
	flag         = "_"
	publishTopic = "$oc/devices/" +
		deviceId +
		"/sys/properties/report"
	serviceId = "deviceInfo"
)

type Message struct {
   services []Services `json:"services"`
}

type Services struct {
	Service_id string     `json:"service_id"`
	Properties Properties `json:"properties"`
	EventTime  string     `json:"event_time"`
}
type Properties struct {
	Temperature float32 `json:"temp"`
	Humidity    float32 `json:"humi"`
}

func main() {
	var wg sync.WaitGroup
	wg.Add(2)
	ch := make(chan int)
	sha, yyyymmddhh := getSha256()
	clientId := deviceId + flag + "0" + flag + "1" + flag + yyyymmddhh
	fmt.Println(clientId)
	ops := mqtt.NewClientOptions().SetClientID(clientId).SetUsername(deviceId).AddBroker(Host).SetPassword(sha).
		SetTLSConfig(getCert()).SetKeepAlive(60 * time.Second).SetCleanSession(false)
	ops.ProtocolVersion = 4
	client := mqtt.NewClient(ops)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())

	} else {
		fmt.Println("CONNECTED!")
	}
	go func() {
		defer wg.Done()
		for {
			<-ch
			log.Print("Publish:")
			pop(client)

		}
	}()

	go func() {
		defer wg.Done()
		for i := 0; i < 10; i++ {
			ch <- 1
			time.Sleep(10 * time.Second)
		}
	}()
	wg.Wait()
	log.Println("Closing")

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

var (
	te float32 = 10.0
	hu float32 = 0.0
)

func mashalMessage() interface{} {
	p := &Properties{
		Temperature: te,
		Humidity:    hu,
	}
	te, hu = te+10.0, hu+10.0

	s := &Services{
		Service_id: serviceId,
		Properties: *p,
		EventTime:  getEventTime(),
	}
	m := new(Message)
	m.Services = []Services{*s}
	mashalMessage, err := json.Marshal(m)
	if err != nil {
		fmt.Println(err)

	}
	fmt.Println(string(mashalMessage))
	return mashalMessage
}

func getEventTime() string {
	var evenTime string
	format := "2006-01-02 15:04:05"
	utc := time.Now().UTC().Format(format)
	utc = strings.Replace(utc, "-", "", -1)
	utc = strings.Replace(utc, " ", "", -1)
	utc = strings.Replace(utc, ":", "", -1)
	for i, v := range utc {
		if i < 8 {
			evenTime += string(v)
		}

	}
	evenTime += "T"
	for i, v := range utc {
		if i >= 8 && i < 14 {
			evenTime += string(v)
		}

	}
	evenTime += "Z"
	return evenTime
}

func pop(client mqtt.Client) mqtt.Client {
	if token := client.Publish(publishTopic, 0, true, mashalMessage());
		token.Wait() && token.Error() != nil {
		panic(token.Error())
	} else {
		fmt.Println("PUBLISHED")
	}
	return client

}
func getCert() *tls.Config {
	tlsconfig := &tls.Config{InsecureSkipVerify: true}
	return tlsconfig

}
