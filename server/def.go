package server

import (
	"flag"
	"github.com/gorilla/websocket"
	"gopkg.in/ini.v1"
	"log"
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
	flag.Parse()

	cfg, err := ini.Load(*c)
	if err != nil {
		panic(err)
	}

	cnf := &config{}
	cnf.Ip = cfg.Section("server").Key("ip").String()
	cnf.PidFile = cfg.Section("").Key("pid_file").String()
	cnf.Port,err = cfg.Section("server").Key("port").Uint64()
	if err != nil {
		panic(err)
	}

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
		c.hub.unregister <- c
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
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := c.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			w.Write(message)
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
