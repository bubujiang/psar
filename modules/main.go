package modules

import (
	"bufio"
	"io"
	"net/http"
	"os"
	"strings"
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

func (p *Pack) SetType(t string)  {
	p.Type = t
}

func (p *Pack) Run(addc func(*Pack) bool)  {
	//需要不断的获得最新的数据
	for{
		path := (*p.Module).FilePath()
		i := strings.Index(path,"http://")
		var fi,err interface{}
		if i==-1{
			fi, err = os.Open((*p.Module).FilePath())
			if err != nil {
				//todo 错误处理
				return
			}
			fi = fi.(*os.File)
		} else {
			fi, err = http.Get((*p.Module).FilePath())
			if err != nil {
				//todo 错误处理
				return
			}
			fi = fi.(*http.Response)
		}

		if fi,ok := fi.(*http.Response);ok{
			err = _read(fi.Body,(*p.Module).Handle)
		}else {
			err = _read(fi,(*p.Module).Handle)
		}

		if err != nil {
			//todo 错误处理
		}

		switch fi.(type) {
		case io.Closer:
			fi.(io.Closer).Close()
			break
		default:
			fi.(*http.Response).Body.Close()
		}
		//todo 写入数据channel(broadcast)
		r := addc(p)
		if !r {
			break
		}
	}
}

func _read(fi interface{}, _handle func(string)) error {
	if fi,ok := fi.(io.Reader);ok {
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
	}else {
		//todo 错误处理
		return nil
	}
}
