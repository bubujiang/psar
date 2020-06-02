package main

import (
	"context"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"strconv"
	"time"
)

type server struct {
	Ip string
	Port uint64
	PidFile string
	wserv *http.Server
}

func (s *server) run(c *config) {
	s.Ip = c.Ip
	s.Port = c.Port
	s.PidFile = c.PidFile
	//s.wserv = &http.Server{}
	s._start()
}

func (s *server) _start()  {
	hub := newHub()
	go hub.run()

	r := gin.Default()
	r.GET("/d", func(c *gin.Context) {
		showData(hub,c)
	})
	//r.Run(s.Ip+":"+strconv.FormatUint(s.Port,10))

	s.wserv = &http.Server{
		Addr:    s.Ip+":"+strconv.FormatUint(s.Port,10),
		Handler: r,
	}
	//s.wserv.Addr = s.Ip+":"+strconv.FormatUint(s.Port,10)
	//s.wserv.Handler = r
	//s.wserv.ListenAndServe()

	if err := s.wserv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("listen: %s\n", err)
	}

	//go func() {
	//	if err := s.wserv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
	//		log.Fatalf("listen: %s\n", err)
	//	}
	//}()
}

func (s *server) stop() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	//todo 所有goroutine退出,清空所有相关channel
	if err := s.wserv.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}

	log.Println("Server exiting")
}

func (s *server) reload() {
	s.stop()
	s._start()
}
