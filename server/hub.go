package server

import (
	"encoding/json"
	"fmt"
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
	fmt.Println(">>>hub run")
	go h.dispatch()
	//var startOnce sync.Once
	isRun := make(chan int)
	for {
		select {
		case client := <-h.register:
			fmt.Println(">>>注册客户端")
			//todo 判断是否启动,没启动就启动监测
			if l := len(h.clients); l<=0 {
				fmt.Println(">>>启动监测")
				isRun = make(chan int)
				for _,pack := range modules.Dpack {
					fmt.Println(">>>下一个模块1")
					go pack.Run(func(p *modules.Pack) bool {
						//fmt.Println(">>>启动监测"+p.Type)
						select {
						case <-isRun:
							return false
						default:
							cp := *p
							d,_ := json.Marshal(cp)
							//todo 错误处理
							h.broadcast <- d
							//fmt.Println(">>>写入公告")
							return true
						}
					})
					fmt.Println(">>>下一个模块2")
				}
				fmt.Println(">>>注册客户端完毕3")
			}
			fmt.Println(">>>注册客户端完毕1")
			h.clients[client] = true
			fmt.Println(">>>注册客户端完毕2")

		case client := <-h.unregister:
			fmt.Println(">>>注销客户端")
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.send)
			}
			//todo 检测客户端数量,关闭监测.
			if l := len(h.clients); l<=0 {
				fmt.Println(">>>关闭监测")
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
				time.Sleep(2)
			default:
				close(client.send)
				delete(h.clients, client)
			}
		}
	}
}
