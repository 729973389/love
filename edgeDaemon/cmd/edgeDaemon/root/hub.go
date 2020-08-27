package root

type Hub struct {
	Down       chan []byte
	Up         chan []byte
	Command    chan []byte
	Register   chan *Client
	UnRegister chan *Client
}

func NewHub() *Hub {
	return &Hub{Down: make(chan []byte, 1024), Up: make(chan []byte, 1024), Register: make(chan *Client), UnRegister: make(chan *Client), Command: make(chan []byte, 1024)}
}

func (hub *Hub) Run() {
	for {
		select {}

	}

}
