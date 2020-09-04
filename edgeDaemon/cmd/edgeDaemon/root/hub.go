package root

import "github.com/wuff1996/edgeDaemon/internal/protobuf"

type Hub struct {
	Down       chan *protobuf.DeviceGister
	Up         chan []byte
	Command    chan []byte
	Register   chan *Client
	UnRegister chan *Client
	DeviceMap  chan []string
}

func NewHub() *Hub {
	return &Hub{DeviceMap: make(chan []string), Down: make(chan *protobuf.DeviceGister, 1024), Up: make(chan []byte, 1024), Register: make(chan *Client), UnRegister: make(chan *Client), Command: make(chan []byte, 1024)}
}

func (hub *Hub) Run() {
	for {
		select {}

	}

}
