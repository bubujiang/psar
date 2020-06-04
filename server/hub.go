package server

import (
	"encoding/json"
	"psar/modules"
)

// Hub maintains the set of active clients and broadcasts messages to the
// clients.
type Hub struct {
	// Registered clients.
	clients map[*Client]bool

	// 写入监测数据
	broadcast chan []byte

	// Register requests from the clients.
	register chan *Client

	// Unregister requests from clients.
	unregister chan *Client
}

func newHub() *Hub {
	return &Hub{
		broadcast:  make(chan []byte,1000),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		clients:    make(map[*Client]bool),
	}
}

func (h *Hub) run() {
	go h.dispatch()
	//var startOnce sync.Once
	isRun := make(chan int)
	for {
		select {
		case client := <-h.register:
			//todo 判断是否启动,没启动就启动监测
			if l := len(h.clients); l<=0 {
				isRun = make(chan int)
				for pack := range modules.Dpack {
					go pack.Run(func(p *modules.Pack) bool {
						select {
						case <-isRun:
							return false
						default:
							cp := *p
							d,_ := json.Marshal(cp)
							//todo 错误处理
							h.broadcast <- d
							return true
						}
					})
				}
			}
			h.clients[client] = true

		case client := <-h.unregister:
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.send)
			}
			//todo 检测客户端数量,关闭监测.
			if l := len(h.clients); l<=0 {
				close(isRun)
				for {
					select {
					case <-h.broadcast:
					default:
						break
					}
				}
			}
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
	for message := range h.broadcast {
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
