package root

import "github.com/wuff1996/edgeDaemon/cmd/edgeDaemon/mqtter"

type Hub struct {
	Down       chan []byte
	Up         chan []byte
	Register   chan *mqtter.Client
	UnRegister chan *mqtter.Client
}

func NewHub() *Hub {
	return &Hub{Down: make(chan []byte, 1024), Up: make(chan []byte, 1024), Register: make(chan *mqtter.Client), UnRegister: make(chan *mqtter.Client)}
}

func (hub *Hub) Run() {
	for {
		select {}

	}

}
