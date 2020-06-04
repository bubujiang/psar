package modules

import (
	"bufio"
	//"encoding/json"
	"io"
	"os"
	//"psar/server"
)

type Module interface {
	Handle(string)
	FilePath() string
	//TimeGap() time.Duration
	Type() string
}

type Pack struct {
	//data *chan interface{}
	Type string
	Module *Module
}

var Dpack []*Pack

func (p *Pack)SetModule(m *Module)  {
	p.Module = m
}

//func (p *Pack) SetData(d *chan interface{}) {
//	p.data = d
//}

func (p *Pack) SetType(t string)  {
	p.Type = t
}

func (p *Pack) Run(addc func(*Pack) bool)  {
	for{
		fi, err := os.Open((*p.Module).FilePath())
		if err != nil {
			//todo 错误处理
			//return err
		}

		err = _read(fi,(*p.Module).Handle)
		if err != nil {
			//todo 错误处理
		}
		fi.Close()
		//todo 写入数据channel(broadcast)
		r := addc(p)
		if !r {
			break
		}
		//cp := *p
		//d,_ := json.Marshal(cp)
		////todo 错误处理
		//server.Thub.broadcast <- d
		//h.broadcast <- d
	}
}

func _read(fi *os.File, _handle func(string)) error {

	br := bufio.NewReader(fi)
	for {
		line, err := br.ReadString('\n')
		if err != nil{
			if err == io.EOF{
				break
			}
			//todo 错误处理
			return err
		}

		_handle(line)
	}

	return nil
}

//func Run(){
//	for pack := range Dpack {
//		go pack.Run()
//	}
//}
