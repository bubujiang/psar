package server

import (
	"encoding/json"
	"psar/modules"
	"sync"
)

// Hub maintains the set of active clients and broadcasts messages to the
// clients.
type Hub struct {
	// Registered clients.
	clients map[*Client]bool

	// 写入监测数据
	Broadcast chan []byte

	// Register requests from the clients.
	register chan *Client

	// Unregister requests from clients.
	unregister chan *Client
}

func newHub() *Hub {
	return &Hub{
		Broadcast:  make(chan []byte,1000),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		clients:    make(map[*Client]bool),
	}
}

func (h *Hub) run() {
	go h.dispatch()
	var startOnce sync.Once
	for {
		select {
		case client := <-h.register:
			h.clients[client] = true
			//todo 启动监测
			startOnce.Do(func() {

				for pack := range modules.Dpack {
					go pack.Run(func(p *modules.Pack) {
						cp := *p
						d,_ := json.Marshal(cp)
						//todo 错误处理
						h.Broadcast <- d
					})
				}

				//modules.Run()
			})
		case client := <-h.unregister:
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.send)
			}
			//todo 检测客户端数量,关闭监测.
		//case message := <-h.broadcast:
		//	for client := range h.clients {
		//		select {
		//		case client.send <- message:
		//		default:
		//			close(client.send)
		//			delete(h.clients, client)
		//		}
		//	}
		}
	}
}

func (h *Hub) dispatch() {
	for message := range h.Broadcast {
		for client := range h.clients {
			select {
			case client.send <- message:
			default:
				close(client.send)
				delete(h.clients, client)
			}
		}
	}
}
