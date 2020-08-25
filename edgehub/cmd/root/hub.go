package root

import (
	"context"
	log "github.com/sirupsen/logrus"
)

//control all message that transformed by each instance
type Hub struct {
	Clients        map[string]*Client
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
		Clients:        make(map[string]*Client),
		Register:       make(chan *Client),
		UnRegister:     make(chan *Client),
		Broadcast:      make(chan []byte, 256),
		HttpMessage:    make(chan []byte, 256),
		HttpUnRegister: make(chan *Client),
		HttpRegister:   make(chan *Client),
	}
}

//select from channel to register and unregister
func (hub *Hub) Run(ctx context.Context) {
	defer log.Warning("Close", "Run")
	for {
		select {
		case client := <-hub.Register:
			if i, ok := hub.Clients[client.SerialNumber]; ok {
				i.Conn.Close()
				delete(hub.Clients, i.SerialNumber)
			}
			log.Println("register:", client.Conn.RemoteAddr())
			hub.Clients[client.SerialNumber] = client
		case client := <-hub.UnRegister:
			if _, ok := hub.Clients[client.SerialNumber]; ok {
				delete(hub.Clients, client.SerialNumber)
				close(client.Send)
				hub.HttpUnRegister <- client
				log.Warning("unregister: ", client.Conn.RemoteAddr())
			}
		case <-ctx.Done():
			close(hub.UnRegister)
			close(hub.Register)
			close(hub.HttpRegister)
			close(hub.HttpMessage)
			close(hub.Broadcast)
			log.Warning("close Hub.Run")
			return
		}

	}
}
