package root

import (
	"log"
	"net"
    protobuf "github.com/wuff1996/edgeDaemon/internal/protobuf"
)

const network string = "tcp"
var tcpAddr,_ = GetAddr()



func init(){
	conn,err := net.DialTCP(network,nil,tcpAddr)
	if err != nil {
		log.Println(err)
	}
	_,err = conn.Write(protobuf.Getbuf(&protobuf.Connect{Id: "wuff",Password: "1996"}))
	if err != nil {
		log.Println(err)

	}
}
func GetAddr()(*net.TCPAddr,error) {
	tcpAddr,err := net.ResolveTCPAddr(network,"192.168.32.175")
	if err!=nil {
		log.Println(err)
		return nil,err
	}
	return tcpAddr,nil
}
