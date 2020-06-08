package main

import (
	"os"
	"os/signal"
	"psar/server"
	"syscall"
)

func main() {
	serv := &server.Serv{}
	serv.Run(server.GetCnf())
	watch(serv)
}

func watch(s *server.Serv) {
	sigs := make(chan os.Signal)
	signal.Notify(sigs,syscall.SIGINT,syscall.SIGUSR1)
	select {
	case sig := <- sigs:
		if sig == syscall.SIGINT{
			s.Stop()
		}else {
			s.Reload()
		}
	}
}
