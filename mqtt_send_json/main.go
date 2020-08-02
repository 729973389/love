package main

import (
	"crypto/hmac"
	"crypto/sha256"
	"crypto/tls"
	"crypto/x509"
	"encoding/hex"
	"encoding/json"
	"fmt"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"math/rand"
	"strings"
	"time"
)

const (
	deviceId = "1755627c-5b04-486c-8d26-124e53dccdaf"
	Host = "tls://" +
		"www.iotplatform-demo.com:8883"
	data = "644a65b32b9e49e0ea84"
	codecMode = "json"
	subscribeTopic =
		"/huawei" +
			"/v1" +
			"/devices" +
			"/" +
		deviceId +
		"/command" +
			"/" +
			codecMode
	publicTopic =
		"/huawei/" +
		"v1/" +
		"devices/" +
		deviceId +
		"/data/" +
		codecMode
	msgType = "deviceReq"

	serviceId = "Test"

)
const mqtt_root ="-----BEGIN CERTIFICATE-----\nMIID4DCCAsigAwIBAgIJAK97nNS67HRvMA0GCSqGSIb3DQEBCwUAMFMxCzAJBgNV\nBAYTAkNOMQswCQYDVQQIEwJHRDELMAkGA1UEBxMCU1oxDzANBgNVBAoTBkh1YXdl\naTELMAkGA1UECxMCQ04xDDAKBgNVBAMTA0lPVDAeFw0xNjA1MDQxMjE3MjdaFw0y\nNjA1MDIxMjE3MjdaMFMxCzAJBgNVBAYTAkNOMQswCQYDVQQIEwJHRDELMAkGA1UE\nBxMCU1oxDzANBgNVBAoTBkh1YXdlaTELMAkGA1UECxMCQ04xDDAKBgNVBAMTA0lP\nVDCCASIwDQYJKoZIhvcNAQEBBQADggEPADCCAQoCggEBAJxM9fwkwvxeILpkvoAM\nGdqq3x0G9o445F6Shg3I0xmmzu9Of8wYuW3c4jtQ/6zscuIGyWf06ke1z//AVZ/o\ndp8LkuFbBbDXR5swjUJ6z15b6yaYH614Ty/d6DrCM+RaU+FWmxmOon9W/VELu2BB\nNXDQHJBSbWrLNGnZA2erk4JSMp7RhHrZ0QaNtT4HhIczFYtQ2lYF+sQJpQMrjoRn\ndSV9WB872Ja4DgcISU1+wuWLmS/NKjIvOWW1upS79yu2I4Rxos2mFy9xxz311rGC\nZ3X65ejFNzCUrNgf6NEP1N7wB9hUu7u50aA+/56D7EgjeI0gpFytC+a4f6JCPVWI\nLr0CAwEAAaOBtjCBszAdBgNVHQ4EFgQUcGqy59oawLEgMl21//7F5RyABpwwgYMG\nA1UdIwR8MHqAFHBqsufaGsCxIDJdtf/+xeUcgAacoVekVTBTMQswCQYDVQQGEwJD\nTjELMAkGA1UECBMCR0QxCzAJBgNVBAcTAlNaMQ8wDQYDVQQKEwZIdWF3ZWkxCzAJ\nBgNVBAsTAkNOMQwwCgYDVQQDEwNJT1SCCQCve5zUuux0bzAMBgNVHRMEBTADAQH/\nMA0GCSqGSIb3DQEBCwUAA4IBAQBgv2PQn66gRMbGJMSYS48GIFqpCo783TUTePNS\ntV8G1MIiQCpYNdk2wNw/iFjoLRkdx4va6jgceht5iX6SdjpoQF7y5qVDVrScQmsP\nU95IFcOkZJCNtOpUXdT+a3N+NlpxiScyIOtSrQnDFixWMCJQwEfg8j74qO96UvDA\nFuTCocOouER3ZZjQ8MEsMMquNEvMHJkMRX11L5Rxo1pc6J/EMWW5scK2rC0Hg91a\nLod6aezh2K7KleC0V5ZlIuEvFoBc7bCwcBSAKA3BnQveJ8nEu9pbuBsVAjHOroVb\n8/bL5retJigmAN2GIyFv39TFXIySw+lW0wlp+iSPxO9s9J+t\n-----END CERTIFICATE-----\n"

type Message struct {
	 MsgType string `json:"msgType"`
	 Date []Data `json:"data"`
}

type Data struct {
	ServiceId   string `json:"serviceId"`
	ServiceData ServiceData `json:"serviceData"`
	EventTime string `json:"eventTime"`
}

type ServiceData struct {
	Temperature int `json:"Temperature"`
	Humidity int `json:"Humidity"`

}


func main(){

	certpool := x509.NewCertPool()
	ok := certpool.AppendCertsFromPEM([]byte(mqtt_root))
	if !ok {
		panic("Failed to parse root certificate")
	}
	tlsconfig := &tls.Config{RootCAs: certpool}

	Password,yyyymmddhh := getSha()
	flag := "_"
	cilentId := deviceId+flag+"0"+flag+"1"+flag+yyyymmddhh
	UserName := deviceId
    ops := mqtt.NewClientOptions().AddBroker(Host).
    	SetClientID(cilentId).
    	SetUsername(UserName).
    	SetPassword(Password).
    	SetKeepAlive(60 * time.Second).
    	SetCleanSession(true).
    	SetTLSConfig(tlsconfig)
    	ops.ProtocolVersion=4
    client := mqtt.NewClient(ops)
    if token := client.Connect() ; token.Wait()&&token.Error()!=nil {
    	panic(token.Error())
	}else {
		fmt.Println("Connect success!")
	}


for i := 0 ;i<10 ;i++ {
	if token := client.Publish(publicTopic,0,true,getPayload()); token.Wait() && token.Error() != nil {
		fmt.Println(token.Error())
	}else {
		token.Wait()
		fmt.Println("Publish sucess :"+publicTopic)
	}
}


}


func getSha () (s string ,d string){
	var yyyymmddhh string

    format := "2006-01-02 15:04:05"
	 utc :=time.Now().UTC().Format(format)
	utc = strings.Replace(utc," ","",-1)
	utc = strings.Replace(utc,"-","",-1)
	for i,v :=range utc {

		if i<10  {
			yyyymmddhh += string(v)

		}

	}

	fmt.Println(data,"__yyyymmddhh:",yyyymmddhh)
	h := hmac.New(sha256.New,[]byte(yyyymmddhh))
	h.Write([]byte(data))
	sha := hex.EncodeToString(h.Sum(nil))
	fmt.Println("shaCode:____ : "+sha)
	return sha,yyyymmddhh
}

func getPayload()interface{}{
	serviceId := new(ServiceData)
	serviceId.setTemperature(rand.Intn(40))
	serviceId.setHumidity(rand.Intn(100))
	data := &Data{
		ServiceId: serviceId,
		ServiceData: *serviceId,
		EventTime: getEventTime(),
	}
	msg := new(Message)
	msg.MsgType = msgType
	msg.Date = []Data{*data}

	me,err := json.Marshal(msg)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(string(me))
	return string(me)

}

func getEventTime() string{
	var yyyymmddhhmmssz string

	format := "2006-01-02 15:04:05"
	yyyymmddhhmmss :=time.Now().UTC().Format(format)
	yyyymmddhhmmss = strings.Replace(yyyymmddhhmmss," ","",-1)
	yyyymmddhhmmss = strings.Replace(yyyymmddhhmmss,"-","",-1)
	yyyymmddhhmmss = strings.Replace(yyyymmddhhmmss,":","",-1)
	yyyymmddhhmmssz = yyyymmddhhmmss+"Z"
	return yyyymmddhhmmssz

}

func(serviceData *ServiceData) setTemperature (i int){
	serviceData.Temperature = i
}

func (serviceData *ServiceData) setHumidity (i int) {
	serviceData.Humidity = i
}

//var msgHandler mqtt.MessageHandler = func(client mqtt.Client, message mqtt.Message) {
//
//
//}











