package server

import (
	"encoding/json"
	"log"
	"psar/modules"
	"time"
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
		broadcast:  make(chan []byte),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		clients:    make(map[*Client]bool),
	}
}

func (h *Hub) run() {
	go h.dispatch()
	isRun := make(chan int)
	for {
		select {
		case client := <-h.register:
			if l := len(h.clients); l<=0 {
				isRun = make(chan int)
				for _,pack := range modules.Dpack {
					go pack.Run(func(p *modules.Pack) bool {
						select {
						case <-isRun:
							return false
						default:
							cp := *p
							d,err := json.Marshal(cp)
							if err != nil {
								log.Printf("error: %v", err)
								return true
							}
							h.broadcast <- d
							time.Sleep(500 * time.Millisecond)
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
			if l := len(h.clients); l<=0 {
				close(isRun)
				for {
					select {
					case <-h.broadcast:
					default:
						goto forEnd
					}
				}
				forEnd:
			}
		}
	}
}

func (h *Hub) dispatch() {
	for message := range h.broadcast {
		for client := range h.clients {
			select {
			case client.send <- message:
			default:
				//close(client.send)
				//delete(h.clients, client)
			}
		}
	}
}
