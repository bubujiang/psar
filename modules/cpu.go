package modules

import (
	//"psar/modules"
	"strconv"
	"strings"
)

const (
	CPUPATH = "proc/stat"
	)

type StatsCpu struct {
	CpuUser uint64
	CpuNice uint64
	CpuSys uint64
	CpuIdle uint64
	CpuIowait uint64
	CpuSteal uint64
	CpuHardirq uint64
	CpuSoftirq uint64
	CpuGuest uint64
	//CpuNumber uint64
}

func (m *StatsCpu)Handle(line string) {
	fields := strings.Fields(line)
	if fields[0]!="cpu"{
		return
	}
	for i,d := range fields[1:] {
		kd,err := strconv.ParseUint(d,10,64)
		if err != nil {
			//todo 错误处理
			return
			//break
		}
		switch i {
		case 0:
			m.CpuUser = kd
		case 1:
			m.CpuNice = kd
		case 2:
			m.CpuSys = kd
		case 3:
			m.CpuIdle = kd
		case 4:
			m.CpuIowait = kd
		case 5:
			m.CpuHardirq = kd
		case 6:
			m.CpuSoftirq = kd
		case 7:
			m.CpuSteal = kd
		case 8:
			m.CpuGuest = kd

		}
	}
	return
}

func (m *StatsCpu)FilePath() string {
	return CPUPATH
}

func (m *StatsCpu)Type() string {
	return "cpu"
}

func init() {
	p := &Pack{}
	cpu := Module(&StatsCpu{})
	p.SetModule(&cpu)
	p.SetType(cpu.Type())
	Dpack = append(Dpack, p)
}