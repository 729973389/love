package root

import (
	"bufio"
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"github.com/wuff1996/edgeHub/config"
	"io/ioutil"
	"os"
)

//filename of configuration
var socketFile = "socket.json"

//Info holds the global configuration
var Info config.Info

//set default configuration
func SetConfig() {
	configSocket := &config.Info{}
	configSocket.EdgeInfoServer = "http://192.168.32.150:8082"
	configSocket.Socket = "43211"
	configSocket.PostEdge = "/api/v1/edge/data/create"
	configSocket.PutStatus = "/api/v1/edge/update/property"
	configSocket.GetInfo = "/api/v1/edge/getInfo"
	configSocket.Key = "3141592666"
	configSocket.PostDevice = "/api/v1/iot/data/transfer"
	configSocket.DeviceInfoServer = "http://192.168.32.11:9000"
	configSocket.DeviceControlServer = "ws://192.168.32.150:8082/api/v1/ws/easyfetch"
	configSocket.GetEdgeDevice = "http://192.168.32.150:8083/api/v1/device/edge/getEdgeDevice"
	configSocket.GetEdgeDeviceToken = "a4b57c3a2241e9cfcf66e415e853d3b2"
	configSocket.EdgeInfoServerToken = "1d38151c19708bd85777b485c3fa9da8"
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

//initialize the configuration
func init() {
	for i := 0; i < 2; i++ {
		b, err := ioutil.ReadFile(socketFile)
		if err != nil {
			if i == 1 {
				panic(fmt.Sprintln("getConfig failed: ", err))
			}
			log.Warning(errors.Wrap(err, "read config"))
			SetConfig()
			continue
		}
		socket := &config.Info{}
		err = json.Unmarshal(b, socket)
		if err != nil {
			log.Error(err)
			panic(fmt.Sprintln("getConfig failed: ", err))
		}
		log.Println("Create default config file success")
		Info = *socket
	}
}
