package main

import (
	"crypto/hmac"
	"crypto/sha256"
	"crypto/tls"
	"crypto/x509"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"strings"
	"sync"
	"time"
)

const (
	deviceId = "1755627c-5b04-486c-8d26-124e53dccdaf"
	Host = "tls://" +
		"www.iotplatform-demo.com:8883"
	secrect = "644a65b32b9e49e0ea84"

)
const mqtt_root ="-----BEGIN CERTIFICATE-----\nMIID4DCCAsigAwIBAgIJAK97nNS67HRvMA0GCSqGSIb3DQEBCwUAMFMxCzAJBgNV\nBAYTAkNOMQswCQYDVQQIEwJHRDELMAkGA1UEBxMCU1oxDzANBgNVBAoTBkh1YXdl\naTELMAkGA1UECxMCQ04xDDAKBgNVBAMTA0lPVDAeFw0xNjA1MDQxMjE3MjdaFw0y\nNjA1MDIxMjE3MjdaMFMxCzAJBgNVBAYTAkNOMQswCQYDVQQIEwJHRDELMAkGA1UE\nBxMCU1oxDzANBgNVBAoTBkh1YXdlaTELMAkGA1UECxMCQ04xDDAKBgNVBAMTA0lP\nVDCCASIwDQYJKoZIhvcNAQEBBQADggEPADCCAQoCggEBAJxM9fwkwvxeILpkvoAM\nGdqq3x0G9o445F6Shg3I0xmmzu9Of8wYuW3c4jtQ/6zscuIGyWf06ke1z//AVZ/o\ndp8LkuFbBbDXR5swjUJ6z15b6yaYH614Ty/d6DrCM+RaU+FWmxmOon9W/VELu2BB\nNXDQHJBSbWrLNGnZA2erk4JSMp7RhHrZ0QaNtT4HhIczFYtQ2lYF+sQJpQMrjoRn\ndSV9WB872Ja4DgcISU1+wuWLmS/NKjIvOWW1upS79yu2I4Rxos2mFy9xxz311rGC\nZ3X65ejFNzCUrNgf6NEP1N7wB9hUu7u50aA+/56D7EgjeI0gpFytC+a4f6JCPVWI\nLr0CAwEAAaOBtjCBszAdBgNVHQ4EFgQUcGqy59oawLEgMl21//7F5RyABpwwgYMG\nA1UdIwR8MHqAFHBqsufaGsCxIDJdtf/+xeUcgAacoVekVTBTMQswCQYDVQQGEwJD\nTjELMAkGA1UECBMCR0QxCzAJBgNVBAcTAlNaMQ8wDQYDVQQKEwZIdWF3ZWkxCzAJ\nBgNVBAsTAkNOMQwwCgYDVQQDEwNJT1SCCQCve5zUuux0bzAMBgNVHRMEBTADAQH/\nMA0GCSqGSIb3DQEBCwUAA4IBAQBgv2PQn66gRMbGJMSYS48GIFqpCo783TUTePNS\ntV8G1MIiQCpYNdk2wNw/iFjoLRkdx4va6jgceht5iX6SdjpoQF7y5qVDVrScQmsP\nU95IFcOkZJCNtOpUXdT+a3N+NlpxiScyIOtSrQnDFixWMCJQwEfg8j74qO96UvDA\nFuTCocOouER3ZZjQ8MEsMMquNEvMHJkMRX11L5Rxo1pc6J/EMWW5scK2rC0Hg91a\nLod6aezh2K7KleC0V5ZlIuEvFoBc7bCwcBSAKA3BnQveJ8nEu9pbuBsVAjHOroVb\n8/bL5retJigmAN2GIyFv39TFXIySw+lW0wlp+iSPxO9s9J+t\n-----END CERTIFICATE-----\n"
func main(){

	certpool := x509.NewCertPool()
	ok := certpool.AppendCertsFromPEM([]byte(mqtt_root))
	if !ok {
		panic("Failed to parse root certificate")
	}
	tlsconfig := &tls.Config{RootCAs: certpool}

	Password,yyyymmddhh := getSha(secrect)
	flag := "_"
	cilentId := deviceId+flag+"0"+flag+"1"+flag+yyyymmddhh
	fmt.Println(cilentId)
	UserName := deviceId
    ops := mqtt.NewClientOptions().AddBroker(Host).
    	SetClientID(cilentId).
    	SetUsername(UserName).
    	SetPassword(Password).
    	SetKeepAlive(60 * time.Second).
    	SetOnConnectHandler( func(client mqtt.Client){fmt.Println("coneceted")}).
    	SetConnectionLostHandler(func(client mqtt.Client,err error){fmt.Println("lost")}).
    	SetCleanSession(true).
    	SetTLSConfig(tlsconfig)
    	ops.ProtocolVersion=4
    client := mqtt.NewClient(ops)
    if token := client.Connect() ; token.Wait()&&token.Error()!=nil {
    	panic(token.Error())
	}else {
		fmt.Println("Connect success!")
	}



}


func getSha (secrect string) (s string ,y string){
	var yyyymmddhh string

    format := "2006-01-02 15:04:05"
	local,err   := time.LoadLocation("UTC")
	if err != nil {
		panic(err)
	}
	 utc :=time.Now().In(local).Format(format)
	utc = strings.Replace(utc," ","",-1)
	utc = strings.Replace(utc,"-","",-1)
	for i,v :=range utc {

		if i<10  {
			yyyymmddhh += string(v)

		}

	}

	fmt.Println(secrect,"__",yyyymmddhh)
	h := hmac.New(sha256.New,[]byte(yyyymmddhh))
	fmt.Println([]byte(secrect))
	_,err = h.Write([]byte(secrect))
	if err != nil {
		panic(err)
	}
	sha := hex.EncodeToString(h.Sum(nil))
	fmt.Println("result : "+sha)
	return sha,yyyymmddhh
}





type Client struct {                     //定义客户端结构体
	nativeClient mqtt.Client
	clientOptions *mqtt.ClientOptions
	locker *sync.Mutex
	//消息收到之后处理函数
	observer func(c *Client,msg *Message)

}

type Message struct {               //定义消息结构体
	ClientId string `json:"client_id"`
	nodeId string
	deviceId string
	temperature string

}


func (client *Client) GetClientID() string {
	return client.clientOptions.ClientID
}




func (client *Client) ensureConnected() error {
	if !client.nativeClient.IsConnected(){
		client.locker.Lock()
		defer client.locker.Unlock()
	}
	if !client.nativeClient.IsConnected() {
		if token := client.nativeClient.Connect();token.Wait()&&token.Error() != nil {
			return token.Error()
		}

	}
	return nil


}

func (client *Client) Publish (topic string, qos byte ,retained bool,data []byte) error{
	//是否连接
	if err := client.ensureConnected(); err != nil {
		return err
	}
	token := client.nativeClient.Publish(topic,qos,retained,data)
	//publish需要传递topic,qos,retained,和一个数据结构
	if err := token.Error(); err != nil {
		return err
	}

	if ok := token.WaitTimeout(time.Second * 10); !ok {
		//连接超时
		return errors.New("mqtt publish wait timeout")
	}
	{

		return nil

	}

	//一切正常



}

func (client *Client) Subscribe (observer func(c *Client,msg *Message),qos byte,topics ...string) error {
	//订阅
	if len(topics) == 0 {
		return errors.New("topic is zero")
	}
	if observer == nil {
		return errors.New("observer is empty")
	}
	client.observer =observer
	filter := make(map[string]byte)
	for _,topic :=range topics {
		filter[topic]=qos
	}
	client.nativeClient.SubscribeMultiple(filter,client.MessageHandler)
	return nil


}
//处理消息
func (client *Client) MessageHandler(c mqtt.Client, msg mqtt.Message) {
	if client.observer ==nil {
		fmt.Println("not subscribe message observe")
		return
	}
	message,err := decodeMessage(msg.Payload())
	if err != nil {
		fmt.Println("failed to decode message")
		return
	}
	client.observer(client,message)
}

func decodeMessage(payload []byte) (*Message,error){
	message := new(Message)
	decoder := json.NewDecoder(strings.NewReader(string(payload)))
	decoder.UseNumber()
	if err := decoder.Decode(&message);err != nil {
		return nil, err
	}
	return message,nil


}

func (client *Client) unsubscribe (topics ...string) {
	client.observer = nil
	client.nativeClient.Unsubscribe(topics...)
}

