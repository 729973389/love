package root

import (
	"bufio"
	"encoding/json"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	config "github.com/wuff1996/edgeDaemon/config"
	"io/ioutil"
	"os"
)

const configFile = "server.json"

func SetConfig() {
	config := config.Server{}
	config.Url = "localhost:43211/hub"
	b, err := json.MarshalIndent(config, "", " ")
	if err != nil {
		log.WithError(err).WithField(configFile, "Error")
		return
	}
	file, err := os.OpenFile(configFile, os.O_CREATE|os.O_RDWR, 0666)
	if err != nil {
		log.Println(errors.Wrap(err, configFile))
		return
	}
	writer := bufio.NewWriter(file)
	writer.Write(b)
	writer.Flush()
}

func GetConfig() *config.Server {
	for i := 0; i < 2; i++ {
		b, err := ioutil.ReadFile(configFile)
		if err != nil {
			log.WithError(err).WithField("http.json", "Error")
			SetConfig()
			if i == 0 {
				continue
			}
			return nil
		}
		config := &config.Server{}
		err = json.Unmarshal(b, config)
		if err != nil {
			log.WithError(err).WithField("http.json", "Error")
			return nil
		}
		log.Println(config)
		return config
	}
	return nil

}
