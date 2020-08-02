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

const root_ca = "-----BEGIN CERTIFICATE-----\nMIID4DCCAsigAwIBAgIJAK97nNS67HRvMA0GCSqGSIb3DQEBCwUAMFMxCzAJBgNV\nBAYTAkNOMQswCQYDVQQIEwJHRDELMAkGA1UEBxMCU1oxDzANBgNVBAoTBkh1YXdl\naTELMAkGA1UECxMCQ04xDDAKBgNVBAMTA0lPVDAeFw0xNjA1MDQxMjE3MjdaFw0y\nNjA1MDIxMjE3MjdaMFMxCzAJBgNVBAYTAkNOMQswCQYDVQQIEwJHRDELMAkGA1UE\nBxMCU1oxDzANBgNVBAoTBkh1YXdlaTELMAkGA1UECxMCQ04xDDAKBgNVBAMTA0lP\nVDCCASIwDQYJKoZIhvcNAQEBBQADggEPADCCAQoCggEBAJxM9fwkwvxeILpkvoAM\nGdqq3x0G9o445F6Shg3I0xmmzu9Of8wYuW3c4jtQ/6zscuIGyWf06ke1z//AVZ/o\ndp8LkuFbBbDXR5swjUJ6z15b6yaYH614Ty/d6DrCM+RaU+FWmxmOon9W/VELu2BB\nNXDQHJBSbWrLNGnZA2erk4JSMp7RhHrZ0QaNtT4HhIczFYtQ2lYF+sQJpQMrjoRn\ndSV9WB872Ja4DgcISU1+wuWLmS/NKjIvOWW1upS79yu2I4Rxos2mFy9xxz311rGC\nZ3X65ejFNzCUrNgf6NEP1N7wB9hUu7u50aA+/56D7EgjeI0gpFytC+a4f6JCPVWI\nLr0CAwEAAaOBtjCBszAdBgNVHQ4EFgQUcGqy59oawLEgMl21//7F5RyABpwwgYMG\nA1UdIwR8MHqAFHBqsufaGsCxIDJdtf/+xeUcgAacoVekVTBTMQswCQYDVQQGEwJD\nTjELMAkGA1UECBMCR0QxCzAJBgNVBAcTAlNaMQ8wDQYDVQQKEwZIdWF3ZWkxCzAJ\nBgNVBAsTAkNOMQwwCgYDVQQDEwNJT1SCCQCve5zUuux0bzAMBgNVHRMEBTADAQH/\nMA0GCSqGSIb3DQEBCwUAA4IBAQBgv2PQn66gRMbGJMSYS48GIFqpCo783TUTePNS\ntV8G1MIiQCpYNdk2wNw/iFjoLRkdx4va6jgceht5iX6SdjpoQF7y5qVDVrScQmsP\nU95IFcOkZJCNtOpUXdT+a3N+NlpxiScyIOtSrQnDFixWMCJQwEfg8j74qO96UvDA\nFuTCocOouER3ZZjQ8MEsMMquNEvMHJkMRX11L5Rxo1pc6J/EMWW5scK2rC0Hg91a\nLod6aezh2K7KleC0V5ZlIuEvFoBc7bCwcBSAKA3BnQveJ8nEu9pbuBsVAjHOroVb\n8/bL5retJigmAN2GIyFv39TFXIySw+lW0wlp+iSPxO9s9J+t\n-----END CERTIFICATE-----"
const (
	deviceId ="5f02cef604557e15f5bd1a9c_testedge"
	secret = "12345678"
	Host = "tls://iot-mqtts.cn-north-4.myhuaweicloud.com:8883"
	flag = "_"
	publishTopic = "oc/devices/"+
		deviceId+"/"+
		"sys/properties/report"
	serviceId ="Test"

)

type Message struct {
	Services []Services `json:"services"`
}

type Services struct {
	Service_id string `json:"service_id"`
	Properties Properties `json:"properties"`
	EventTime string `json:"event_time"`
}
type Properties struct {
	Temperature int `json:"Temperature"`
	Humidity int `json:"Humidity"`

}




func main(){
	sha,yyyymmddhh := getSha256()
	clientId := deviceId+flag+"0"+flag+"1"+flag+yyyymmddhh
	fmt.Println(clientId)
	ops := mqtt.NewClientOptions().SetClientID(clientId).SetUsername(deviceId).AddBroker(Host).SetPassword(sha).
		SetTLSConfig(getCert()).SetKeepAlive(60 * time.Second).SetCleanSession(false)
	ops.ProtocolVersion=4
	client := mqtt.NewClient(ops)
	if token := client.Connect(); token.Wait()&&token.Error()!=nil {
		panic(token.Error())

	} else {
		fmt.Println("CONNECTED!")
	}
	pop(client)

}

func getSha256()(s string,t string){
	var yyyymmddhh string
	format := "2006-01-02 15:04:05"
	utc := time.Now().UTC().Format(format)
	utc = strings.Replace(utc,"-","",-1)
	utc = strings.Replace(utc," ","",-1)
	utc = strings.Replace(utc,":","",-1)
	for i,v := range utc{
		if i<10 {
			yyyymmddhh+=string(v)

		}
	}
	fmt.Println("TIME FLAG: "+yyyymmddhh)
	sha := hmac.New(sha256.New,[]byte(yyyymmddhh))
	sha.Write([]byte(secret))
	fmt.Println(hex.EncodeToString(sha.Sum(nil)))
	return hex.EncodeToString(sha.Sum(nil)),yyyymmddhh

}

func mashalMessage () interface{} {
	p :=&Properties{
		Temperature: rand.Intn(40),
		Humidity: rand.Intn(100),
	}

	s := &Services{
		Service_id: serviceId,
		Properties: *p,
		EventTime: getEventTime(),
	}
	m := new(Message)
	m.Services = []Services{*s}
	mashalMessage,err := json.Marshal(m)
	if err!=nil {
		fmt.Println(err)

	}
	fmt.Println(string(mashalMessage))
	return mashalMessage
}



func getEventTime()string{
	var evenTime string
	format := "2006-01-02 15:04:05"
	utc := time.Now().UTC().Format(format)
	utc = strings.Replace(utc,"-","",-1)
	utc = strings.Replace(utc," ","",-1)
	utc = strings.Replace(utc,":","",-1)
	for i,v := range utc{
		if i<8 {
			evenTime+=string(v)
		}

	}
	evenTime+="T"
	for i,v := range utc{
		if i>=8&&i<14 {
			evenTime+=string(v)
		}

	}
	evenTime+="Z"
	return evenTime
}

func  pop(client mqtt.Client) mqtt.Client {
	if token := client.Publish(publishTopic,0,true,mashalMessage());
	token.Wait()&&token.Error()!=nil{
		panic(token.Error())
	} else {
		fmt.Println("PUBLISHED")
	}
	return client

}

func getCert() *tls.Config{
	certpool := x509.NewCertPool()
	certpool.AppendCertsFromPEM([]byte(root_ca))
	tlsconfig := &tls.Config{InsecureSkipVerify: true}
	return tlsconfig

}

