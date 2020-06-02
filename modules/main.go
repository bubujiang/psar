package modules

import (
	"bufio"
	"io"
	"os"
	//_ "psar/modules/mem"
	"time"
)

type Module interface {
	Handle(string)
	FilePath() string
	TimeGap() time.Duration
	Type() string
}

type Pack struct {
	//data *chan interface{}
	Type string
	Module *Module
}

var Dpack  = make(chan *Pack,100)

func (p *Pack)SetModule(m *Module)  {
	p.Module = m
}

//func (p *Pack) SetData(d *chan interface{}) {
//	p.data = d
//}

func (p *Pack) SetType(t string)  {
	p.Type = t
}

func (p *Pack) Run()  {
	for{
		fi, err := os.Open((*p.Module).FilePath())
		if err != nil {
			//todo 错误处理
			//return err
		}

		err = p.read(fi,(*p.Module).Handle)
		if err != nil {
			//todo 错误处理
		}
		fi.Close()
		//
		//ch <- x
		//boot.Data <- *p
		//time.Sleep((*p.Module).TimeGap())
	}
}

func (p *Pack)read(fi *os.File, _handle func(string)) error {

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

func Run(){
	for pack := range Dpack {
		pack.Run()
	}
}
