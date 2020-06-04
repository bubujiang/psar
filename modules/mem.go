package modules

import (
	//"psar/modules"
	"strconv"
	"strings"
)

const (
	TYPE = "mem"
	FILEPATH = "proc/meminfo"
	//TIMEGAP = time.Microsecond * 100
	)

type Stats struct {
	Frmkb uint64
	Bufkb uint64
	Camkb uint64
	Tlmkb uint64
	Acmkb uint64
	Iamkb uint64
	Slmkb uint64
	Frskb uint64
	Tlskb uint64
	Caskb uint64
	Comkb uint64
}

func (m *Stats)Handle(line string) {
	//tag := line[:5]
	fields := strings.Fields(line)
	kb,err := strconv.ParseUint(fields[1],10,64)
	if err != nil {
		//todo 错误处理
		return
		//break
	}
	switch fields[0] {
	case "MemTotal:":
		m.Tlmkb = kb
		break
	case "MemFree:":
		m.Frmkb = kb
		break
	case "Buffers:":
		m.Bufkb = kb
		break
	case "Cached:":
		m.Camkb = kb
		break
	case "Active:":
		m.Acmkb = kb
		break
	case "Inactive:":
		m.Iamkb = kb
		break
	case "Slab:":
		m.Slmkb = kb
		break
	case "SwapCached:":
		m.Caskb = kb
		break
	case "SwapTotal:":
		m.Tlskb = kb
		break
	case "SwapFree:":
		m.Frskb = kb
		break
	case "Committed_AS:":
		m.Comkb = kb
		break
	}

	return
}

func (m *Stats)FilePath() string {
	return FILEPATH
}

//func (m *Stats)TimeGap() time.Duration {
//	return TIMEGAP
//}

func (m *Stats)Type() string {
	return TYPE
}

func init() {
	p := &Pack{}
	mem := Module(&Stats{})
	p.SetModule(&mem)
	p.SetType(mem.Type())
	Dpack = append(Dpack, p)
}