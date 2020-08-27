package root

import (
	"encoding/json"
	log "github.com/sirupsen/logrus"
	config "github.com/wuff1996/edgeDaemon/config"
	"io/ioutil"
)

const configFile = "server.json"

type Server struct {
	Url          string `json:"url"`
	Token        string `json:"token"`
	SerialNumber string `json:"serialNumber"`
}

func GetConfig() *config.Server {
	config := &config.Server{}
	b, err := ioutil.ReadFile(configFile)
	if err != nil {
		log.Fatal("server.json", err)
		return config
	}
	err = json.Unmarshal(b, config)
	if err != nil {
		log.Fatal("server.json", err)
		return config
	}
	log.Println(config)
	return config

}
