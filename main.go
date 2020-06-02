package main

import (
	"os"
	"os/signal"
	"syscall"
)

func main() {
	serv := &server{}
	serv.run(getCnf())
	watch(serv)
}

func watch(s *server) {
	sigs := make(chan os.Signal)
	signal.Notify(sigs,syscall.SIGINT,syscall.SIGUSR1)
	select {
	case sig := <- sigs:
		if sig == syscall.SIGINT{
			s.stop()
		}else {
			s.reload()
		}
	}
}
