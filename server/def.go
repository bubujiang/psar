package server

import (
	"flag"
	"fmt"
	"github.com/gorilla/websocket"
	"gopkg.in/ini.v1"
	"log"
	"os"
	"time"
)

/**
 * 解析配置
 */
type config struct {
	Ip string
	Port uint64
	PidFile string
}

func GetCnf() *config {
	c := flag.String("c","conf.ini","配置文件路径")
	//flag.String("c","conf.ini","配置文件路径")
	flag.Parse()

	cfg, err := ini.Load(*c)
	if err != nil {
		fmt.Printf("Fail to read file: %v", err)
		os.Exit(1)
	}

	cnf := &config{}
	cnf.Ip = cfg.Section("server").Key("ip").String()
	cnf.Port,_ = cfg.Section("server").Key("port").Uint64()
	cnf.PidFile = cfg.Section("").Key("pid_file").String()

	return cnf
}

/**
 * 操作客户端请求
 */
type Client struct {
	conn *websocket.Conn
	send chan []byte
	hub *Hub
}

func (c *Client) readPump() {
	defer func() {
		fmt.Println(">>>断开1")
		c.hub.unregister <- c
		fmt.Println(">>>断开2")
		c.conn.Close()
	}()
	c.conn.SetReadLimit(maxMessageSize)
	c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(func(string) error { c.conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })
	for {
		_, _, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}
			break
		}
		//if x ==
		//message = bytes.TrimSpace(bytes.Replace(message, newline, space, -1))
		//c.hub.broadcast <- message
	}
}

func (c *Client) writePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()
	for {
		select {
		case message, ok := <-c.send:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				// The hub closed the channel.
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := c.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			w.Write(message)

			// Add queued chat messages to the current websocket message.
			//n := len(c.send)
			//for i := 0; i < n; i++ {
			//	w.Write(newline)
			//	w.Write(<-c.send)
			//}

			//if err := w.Close(); err != nil {
			//	return
			//}
		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

/**
 * 监控
 */
