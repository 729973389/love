package root

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/wuff1996/edgeHub/internal/protobuf"
	"log"
	"net/http"
)

var urlFalg = flag.String("url", "192.168.32.150:8081/api/v2/edge/data/create", "set a specific url to connect")

func init(){
	flag.Parse()
}
func Post() {
	s := protobuf.GetSystemInfo()
	interfaceEdge := &protobuf.InterfaceEdge{SerialNumber: "1",Data: &s}
	b,err := json.MarshalIndent(interfaceEdge,""," ")
	if err != nil {
		log.Println("jsonMarshal: ",err)
	}
	fmt.Println(string(b))
	resp,err:= http.Post("http://192.168.32.150:8081/api/v2/edge/data/create","application/json",bytes.NewReader(b))
	if err != nil {
		log.Println("Get: ",err)
		return
	}
	defer resp.Body.Close()
	respBuf:=make([]byte,256)
	resp.Body.Read(respBuf)
	log.Println(string(respBuf))

}
