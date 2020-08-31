package root

import (
	"bufio"
	"encoding/json"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"github.com/wuff1996/edgeHub/config"
	"io/ioutil"
	"os"
)

var socketFile = "socket.json"

func SetConfig() {
	configSocket := &config.Info{}
	configSocket.Url = "http://192.168.32.150:8081"
	configSocket.Socket = "43211"
	configSocket.PostEdge = "/api/v2/edge/data/create"
	configSocket.PutStatus = "/api/v2/edge/update/online"
	configSocket.GetInfo = "/api/v2/edge/getInfo"
	configSocket.Key = "3141592666"
	configSocket.PostDevice = "/api/v1/iot/data/transfer"
	//demo configSocket.GetCommand="xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx"
	bs, err := json.MarshalIndent(configSocket, "", " ")
	if err != nil {
		log.WithError(err).WithField(socketFile, "Error")
		return
	}
	files, err := os.OpenFile(socketFile, os.O_CREATE|os.O_RDWR, 0666)
	if err != nil {
		log.Println(errors.Wrap(err, socketFile))
		return
	}
	writers := bufio.NewWriter(files)
	writers.Write(bs)
	writers.Flush()
}

func GetConfig() config.Info {
	for i := 0; i < 2; i++ {
		b, err := ioutil.ReadFile(socketFile)
		if err != nil {
			log.Warning(errors.Wrap(err, "read config"))
			SetConfig()
			if i == 0 {
				continue
			}
		}
		socket := &config.Info{}
		err = json.Unmarshal(b, socket)
		if err != nil {
			log.Warning(err)
		}
		log.Println("Create default config file success")
		return *socket

	}
	return config.Info{}
}
