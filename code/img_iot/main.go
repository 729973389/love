package main

import (
	"bufio"
	"crypto/hmac"
	"crypto/sha256"
	"crypto/tls"
	"encoding/hex"
	"encoding/json"
	"fmt"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"gocv.io/x/gocv"
	"log"
	"os"
	"strings"
	"time"
)

func main() {
		src  := gocv.NewMat()
		img  := gocv.NewMat()
	src.ConvertTo(&img,gocv.MatTypeCV8SC1)
		defer src.Close()
		defer img.Close()
	vc, err := gocv.OpenVideoCapture("rtsp://admin:C79681197@192.168.32.12:554")
	if err != nil {
		panic(err)
	}
	
	defer vc.Close()

	//frame := vc.Get(gocv.VideoCaptureFrameCount)
	//vc.Set(gocv.VideoCapturePosFrames,frame)

	vc.Read(&img)
	buff, err := gocv.IMEncode(".jpg", img)
	if err != nil {
		panic(err)

	}
	log.Println(len(buff))
	file, err := os.OpenFile(timeToString()+".jpg", os.O_CREATE|os.O_RDWR, 0666)
	if err != nil {
		panic(err)
	}
	defer file.Close()
	newWriter := bufio.NewWriter(file)
	if _, err := newWriter.Write(buff); err != nil {
		panic(err)
	}
	if err = newWriter.Flush(); err != nil {
		panic(err)
	}

	sha, yyyymmddhh := GetSha256()
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
	defer client.Disconnect(6000)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())

	} else {
		fmt.Println("CONNECTED!")
	}
	log.Println(client.IsConnected())
	i := 0
	for ; i < len(buff)-10000; i += 10000 {
		Pop(client, buff[i:i+10000])
		time.Sleep(1 * time.Second)
	}
	Pop(client, buff[i:])

}

func timeToString() string {
	var stringTime string
	format := "2006-01-02 15:04:05"
	t, err := time.LoadLocation("Asia/Shanghai")
	if err != nil {
		log.Print(err)
	}
	l := time.Now().In(t).Format(format)
	l = strings.Replace(l, "-", "", -1)
	l = strings.Replace(l, " ", "", -1)
	l = strings.Replace(l, ":", "", -1)
	for _, v := range l {
		stringTime += string(v)
	}
	return stringTime
}

const (
	deviceId     = "5f0bae9b52045702fceed061_myimg"
	secret       = "12345678"
	Host         = "tls://192.168.32.189:7883"
	flag         = "_"
	publishTopic = "$oc/devices/" + deviceId +
		"/sys/properties/report"
	serviceId = "Image"
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
	Frame []byte `json:"frame"`
}

func GetSha256() (s string, t string) {
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

func MashalMessage(b []byte) interface{} {

	p1 := &Properties{
		Frame: b,
	}

	s := &Services{
		Service_id: serviceId,
		Properties: *p1,
		EventTime:  GetEventTime(),
	}
	m := Message{
		Services: []Services{*s},
	}
	mashalMessage, err := json.Marshal(m)
	if err != nil {
		fmt.Println(err)

	}
	fmt.Println(len(string(mashalMessage)))
	return string(mashalMessage)
}

func GetEventTime() string {
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

func Pop(client mqtt.Client, b []byte) {
	if token := client.Publish(publishTopic, 0, false, MashalMessage(b));
		token.Wait() && token.Error() != nil {
		panic(token.Error())
	} else {
		fmt.Println("PUBLISHED")
	}
}
