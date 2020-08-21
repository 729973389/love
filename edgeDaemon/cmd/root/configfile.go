package root

import (
	"encoding/json"
	log "github.com/sirupsen/logrus"
	config "github.com/wuff1996/edgeDaemon/config"
	"io/ioutil"
)

const configFile = "server.json"

//
//var remoteUrl = flag.String("url", "", "set remote url")
//var installerToken = flag.String("token", "", "set token")
//var number = flag.String("id", "", "set device id")

type Server struct {
	Url          string `json:"url"`
	Token        string `json:"token"`
	SerialNumber string `json:"serialNumber"`
}

//func init() {
//	flag.Parse()
//	if *remoteUrl == "" || *installerToken == "" || *number == "" {
//		log.Fatalf("ERROR: No specific parameters!")
//	}
//	file, err := os.OpenFile("server.json", os.O_CREATE|os.O_RDWR, 0666)
//	defer file.Close()
//	if err != nil {
//		log.Error(err)
//		os.Exit(0)
//	}
//	writer := bufio.NewWriter(file)
//	server := &Server{
//		Url:          *remoteUrl,
//		Token:        *installerToken,
//		SerialNumber: *number,
//	}
//	bs, err := json.MarshalIndent(server, "", " ")
//	if err != nil {
//		log.Error(err)
//	}
//	writer.Write(bs)
//	writer.Flush()
//}

func GetConfig() *config.Server {
	config := &config.Server{}
	b, err := ioutil.ReadFile(configFile)
	if err != nil {
		log.Error("http.json", "Error")
		return config
	}
	err = json.Unmarshal(b, config)
	if err != nil {
		log.Error("http.json", "Error")
		return config
	}
	log.Println(config)
	return config

}
