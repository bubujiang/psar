package main

import "psar/server"

func main() {
	serv := &server.Serv{}
	serv.Run(server.GetCnf())
	//watch(serv)
}

//func watch(s *server) {
//	sigs := make(chan os.Signal)
//	signal.Notify(sigs,syscall.SIGINT,syscall.SIGUSR1)
//	select {
//	case sig := <- sigs:
//		if sig == syscall.SIGINT{
//			s.stop()
//		}else {
//			s.reload()
//		}
//	}
//}
