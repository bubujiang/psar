package server

import (
	"context"
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
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
	s._start()
}

func (s *Serv) _start()  {
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

	go func() {
		if err := s.wserv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			panic(err)
		}
	}()
}

func (s *Serv) Stop() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := s.wserv.Shutdown(ctx); err != nil {
		panic(err)
	}

	os.Exit(0)
}

func (s *Serv) Reload() {
	s.Stop()
	s._start()
}
