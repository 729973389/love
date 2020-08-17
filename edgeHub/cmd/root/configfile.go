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
var socketFile="socket.json"
func SetConfig()  {
	configSocket := &config.Socket{}
	configSocket.Socket = "43211"
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

func GetConfig()config.Socket{
	for i:=0;i<2;i++ {
		b,err := ioutil.ReadFile(socketFile)
		if err != nil{
			log.Warning(errors.Wrap(err,"read config"))
			SetConfig()
			if i==0{
				continue
			}
		}
		socket := &config.Socket{}
		err =json.Unmarshal(b,socket)
		if err != nil {
			log.Warning(err)
		}
		return *socket

	}
	return config.Socket{}
}
