package main

import (
	"bufio"
	"encoding/json"
	"flag"
	log "github.com/sirupsen/logrus"
	"os"
)

var remoteUrl = flag.String("url", "", "set remote url")
var installerToken = flag.String("token", "", "set token")
var number = flag.String("id", "", "set device id")

type Server struct {
	Url          string `json:"url"`
	Token        string `json:"token"`
	SerialNumber string `json:"serialNumber"`
}

func main() {
	flag.Parse()
	if *remoteUrl == "" || *installerToken == "" || *number == "" {
		log.Fatalf("ERROR: No specific parameters!")
	}
	if err := os.Remove("server.json"); err != nil {
	}
	file, err := os.OpenFile("server.json", os.O_CREATE|os.O_RDWR, 0666)
	defer file.Close()
	if err != nil {
		log.Error(err)
		os.Exit(0)
	}
	writer := bufio.NewWriter(file)
	server := &Server{
		Url:          *remoteUrl,
		Token:        *installerToken,
		SerialNumber: *number,
	}
	bs, err := json.MarshalIndent(server, "", " ")
	if err != nil {
		log.Error(err)
	}
	writer.Write(bs)
	writer.Flush()
}
