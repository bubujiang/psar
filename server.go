package main

type server struct {
	Ip string
	Port uint64
	PidFile string
}

func (s *server) run(c *config) {
	s.Ip = c.Ip
	s.Port = c.Port
	s.PidFile = c.PidFile
}

func (s *server) _start()  {

}

func (s *server) stop() {

}

func (s *server) reload() {

}
