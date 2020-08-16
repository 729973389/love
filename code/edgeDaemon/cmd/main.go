package main

import (
	"encoding/json"
	"github.com/wuff1996/edgeDaemon/internal/protobuf"
	"log"
)

func main () {
	 // root.Connet()

	//protobuf.GetCPUInfo()
	s:=protobuf.GetSystemInfo()
	b,_:=json.MarshalIndent(s,""," ")
	log.Printf("%s",b)
}