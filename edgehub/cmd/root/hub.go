package root

import (
	log "github.com/sirupsen/logrus"
)

type Hub struct {
	Clients        map[*Client]bool
	Register       chan *Client
	UnRegister     chan *Client
	Broadcast      chan []byte
	HttpMessage    chan []byte
	HttpUnRegister chan *Client
	HttpRegister   chan *Client
}

//create a hub that control the client lifecycle
func NewHub() *Hub {
	return &Hub{
		Clients:        make(map[*Client]bool),
		Register:       make(chan *Client),
		UnRegister:     make(chan *Client),
		Broadcast:      make(chan []byte, 256),
		HttpMessage:    make(chan []byte, 256),
		HttpUnRegister: make(chan *Client, 256),
		HttpRegister:   make(chan *Client, 256),
	}
}

//select from channel to register and unregister
func (hub *Hub) Run() {
	for {
		select {
		case client := <-hub.Register:
			log.Println("register:", client.Conn.RemoteAddr())
			hub.Clients[client] = true
		case client := <-hub.UnRegister:
			if _, ok := hub.Clients[client]; ok {
				delete(hub.Clients, client)
				close(client.Send)
				hub.HttpUnRegister <- client
				log.Warning("unregister: ", client.Conn.RemoteAddr())
			}

		}

	}
}
