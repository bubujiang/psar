package server

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"strconv"
	"time"
)

type Serv struct {
	Ip string
	Port uint64
	PidFile string
	wserv *http.Server
}

var Thub *Hub

func (s *Serv) Run(c *config) {
	s.Ip = c.Ip
	s.Port = c.Port
	s.PidFile = c.PidFile
	//s.wserv = &http.Server{}
	s._start()
}

func (s *Serv) _start()  {
	fmt.Println(">>>server start")
	Thub = newHub()
	go Thub.run()

	r := gin.Default()
	r.GET("/d", func(c *gin.Context) {
		showData(Thub,c)
	})

	s.wserv = &http.Server{
		Addr:    s.Ip+":"+strconv.FormatUint(s.Port,10),
		Handler: r,
	}

	if err := s.wserv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("listen: %s\n", err)
	}

	//go func() {
	//	if err := s.wserv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
	//		log.Fatalf("listen: %s\n", err)
	//	}
	//}()
}

func (s *Serv) Stop() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	//todo 所有goroutine退出,清空所有相关channel
	if err := s.wserv.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}

	log.Println("Server exiting")
}

func (s *Serv) Reload() {
	s.Stop()
	s._start()
}
