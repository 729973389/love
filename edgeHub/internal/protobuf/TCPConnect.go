//go:generate protoc --proto_path=. --go_out=. TCPConnect.proto

package protobuf

import (
	"github.com/golang/protobuf/proto"
	"log"
)

func GetBuf (connect *Connect) []byte{
	b,err := proto.Marshal(connect)
	if err != nil {
		log.Println(err)
	}
	return b
}

func ReadBuf (b []byte) Connect{
	c := Connect{}
	err := proto.Unmarshal(b,&c)
	if err != nil {
		log.Println(err)
	}
	return c
}