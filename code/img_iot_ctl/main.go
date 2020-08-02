package main

import (
	"bufio"   //bufio.NewWriter(File).Write().Flush()
	"crypto/hmac"  //hmac.New(Sha256.New,yyyymmddhh).Write(data)
	"crypto/sha256"
	"crypto/tls"   //tls.TlsConfig{InsecureSkipVerify:true}
	"encoding/hex"  //hex.EncodeToString( hash.sum(nil))
	"encoding/json"  //json.Marshal(interface{})  json.Unmarshal([]byte,interface{})
	mqtt "github.com/eclipse/paho.mqtt.golang"//mqtt.NewClientOptions().SetBroker().... mqtt.MessageHandler mqtt.Message mqtt.Client
	"gocv.io/x/gocv"  //gocv.NewMat()  gocv.OpenVideoCapture(url string) VideoCapture.Read(*Mat)   gocv.IMEncode(FileExt,Mat)
	"log"   //log.Println()  log.Fatal()
	"os"    //os.OpenFile(name string,os.O_CREAT|O_RDWR,0666 perm)
	"strings"  //strings.Replace(string,old string,new string,-1 int)
	"sync" //wg:=sync.WaitGroup  wg.Add() defer wg.Done() wg.Wait()
	"time"  //time.Now().In(time.LoadLocation()).Format(format string)  time.Sleep(t time.Duration)
)

const (
	deviceId     = "5f0bae9b52045702fceed061_myimg"
	secret       = "12345678"
	Host         = "tls://192.168.32.189:7883"
	flag         = "_"
	publishTopic = "$oc/devices/" + deviceId +
		"/sys/properties/report"
	serviceId      = "Image"
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
	Interval int `json:"interval"`
}

type Response struct {
	ResultCode    int           `json:"result_code"`
	ResponseName  string        `json:"response_name"`
	ResponseParas ResponseParas `json:"paras"`
}
type ResponseParas struct {
	Result int `json:"ack"`
}

var signal = 1
var messageHandler mqtt.MessageHandler = func(client mqtt.Client, message mqtt.Message) {
	if strings.Contains(message.Topic(), "response") {
		return
	}
	commands := new(Commands)
	log.Println("Handling message" + message.Topic())
	err := json.Unmarshal(message.Payload(), commands)
	if err != nil {
		panic(err)
	}
	signal = commands.Paras.Interval
	sendTopic := strings.Replace(message.Topic(), "commands", "commands/response", -1)
	if token := client.Publish(sendTopic, 0, false, getResponseMSG()); token.Wait() && token.Error() != nil {
		log.Println(err)
	} else {
		log.Println("Response success!")
	}
}

func main() {
	var wg sync.WaitGroup
	wg.Add(2)
	ch := make(chan string)
	sha, yyyymmddhh := GetSha256()
	clientId := deviceId + flag + "0" + flag + "1" + flag + yyyymmddhh
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
		log.Println("CONNECTED!------" + clientId + "------")
	}
	go func() {
		defer wg.Done()
		for {
			switch {
			case signal == 0:
				log.Println("*******Stop Capturing*******")
				ch <- "Stop Capturing****"
				signal = -1
				time.Sleep(10 * time.Second)
				continue
			case signal > 0:
				for i := 0; i < signal; i++ {
					time.Sleep(1 * time.Minute)
				}
				getFrameChan(&ch)
				continue
			case signal == -1:
				time.Sleep(20 * time.Second)
				continue
			}
		}
	}()
	go func() {
		defer wg.Done()
		for {
			message := <-ch
			Pop(&client, message)
		}
	}()
	go subscrebe(client)
	wg.Wait()
	log.Println("closing")
	time.Sleep(10 * time.Second)
}

func getFrame() string {
	src := gocv.NewMat()
	img := gocv.NewMat()
	defer src.Close()
	defer img.Close()
	src.ConvertTo(&img, gocv.MatTypeCV8SC1)
	vc, err := gocv.OpenVideoCapture("rtsp://admin:C79681197@192.168.32.12:554")
	if err != nil {
		log.Println(err)
	}
	defer vc.Close()
	//frame := vc.Get(gocv.VideoCaptureFrameCount)
	//vc.Set(gocv.VideoCapturePosFrames,frame)
	vc.Read(&img)
	buff, err := gocv.IMEncode(".jpg", img)
	if err != nil {
		log.Println(err)
	}
	name := timeToString() + ".jpg"
	file, err := os.OpenFile(name, os.O_CREATE|os.O_RDWR, 0666)
	if err != nil {
		log.Println(err)
	}
	defer file.Close()
	newWriter := bufio.NewWriter(file)
	if _, err := newWriter.Write(buff); err != nil {
		log.Println(err)
	}
	if err = newWriter.Flush(); err != nil {
		log.Println(err)
	}
	return name
}

func getFrameChan(ch *chan string) {
	name := getFrame()
	*ch <- name
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

func getResponseMSG() string {
	p := new(ResponseParas)
	p.Result = 1
	r := &Response{
		ResultCode:    0,
		ResponseParas: *p,
		ResponseName:  "CMD",
	}
	b, err := json.Marshal(*r)
	if err != nil {
		log.Println(err)
	}
	return string(b)

}

func subscrebe(c mqtt.Client) {
	token := c.Subscribe(subscribeTopic, 0, messageHandler)
	if token.Wait() && token.Error() != nil {
		panic(token.Error())
	}
	log.Println("Subscribe : " + subscribeTopic + ":OK")
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
	sha := hmac.New(sha256.New, []byte(yyyymmddhh))
	sha.Write([]byte(secret))
	return hex.EncodeToString(sha.Sum(nil)), yyyymmddhh

}

func MarshalMessage(b string) interface{} {

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
	marshalMessage, err := json.Marshal(m)
	if err != nil {
		log.Println(err)

	}
	return string(marshalMessage)
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
	return evenTime
}

func Pop(client *mqtt.Client, b string) {
	if token := (*client).Publish(publishTopic, 1, false, MarshalMessage(b));
		token.Wait() && token.Error() != nil {
		log.Println(token.Error())
	} else {
		log.Println("PUBLISHED")
	}
}
