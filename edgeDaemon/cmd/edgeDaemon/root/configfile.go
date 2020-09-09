package root

import (
	"encoding/json"
	log "github.com/sirupsen/logrus"
	"github.com/wuff1996/edgeDaemon/config"
	"io/ioutil"
)

//filename of configFile
const configFile = "server.json"

//Config holds the global configuration
var Config config.Server

//specific the configuration
func init() {
	c := &config.Server{}
	b, err := ioutil.ReadFile(configFile)
	if err != nil {
		log.Fatal("server.json: ", err)
	}
	err = json.Unmarshal(b, c)
	if err != nil {
		log.Fatal("server.json: ", err)
	}
	Config = *c
}
