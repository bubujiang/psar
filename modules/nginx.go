package modules

import (
	"regexp"
	"strconv"
	"strings"
)

const (
	NGPATH = "http://108.61.85.27:8080/nginx_status"
)

type Nginx struct {
	Naccept uint64
	Nhandled uint64
	Nrequest uint64
	Nactive uint64
	Nreading uint64
	Nwriting uint64
	Nwaiting uint64
	Nrstime uint64
	Nspdy uint64
	Nsslhds uint64
	Nssl uint64
	Nsslk uint64
	Nsslf uint64
	Nsslv3f uint64
	Nhttp2 uint64
}

func (m *Nginx)Handle(line string) {
	//tag := line[:5]
	fields := strings.Fields(line)
	if fields[0]=="Active" && fields[1]=="connections:" {
		m.Nactive,_ = strconv.ParseUint(fields[2],10,64)
	}else if IsDigit,_ := regexp.MatchString("\\d+",fields[0]);IsDigit{
		m.Naccept,_ = strconv.ParseUint(fields[0],10,64)
		m.Nhandled,_ = strconv.ParseUint(fields[1],10,64)
		m.Nrequest,_ = strconv.ParseUint(fields[2],10,64)
	}else if fields[0]=="Reading:" {
		m.Nreading,_ = strconv.ParseUint(fields[1],10,64)
		m.Nwriting,_ = strconv.ParseUint(fields[3],10,64)
		m.Nwaiting,_ = strconv.ParseUint(fields[5],10,64)
	}
	return
}

func (m *Nginx)FilePath() string {
	return NGPATH
}

//func (m *Stats)TimeGap() time.Duration {
//	return TIMEGAP
//}

func (m *Nginx)Type() string {
	return "nginx"
}

func init() {
	p := &Pack{}
	mem := Module(&Nginx{})
	p.SetModule(&mem)
	p.SetType(mem.Type())
	Dpack = append(Dpack, p)
}