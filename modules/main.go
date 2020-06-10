package modules

import (
	"bufio"
	"github.com/pkg/errors"
	"io"
	"log"
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
				log.Printf("error: %v", err)
				return
			}
		} else {
			fi, err = http.Get((*p.Module).FilePath())
			if err != nil {
				log.Printf("error: %v", err)
				return
			}
		}

		switch fi.(type) {
		case io.Closer:
			err = _read(fi,(*p.Module).Handle)
			fi.(io.Closer).Close()
			break
		default:
			err = _read(fi.(*http.Response).Body,(*p.Module).Handle)
			fi.(*http.Response).Body.Close()
		}

		if err != nil {
			log.Printf("error: %v", err)
			return
		}

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
				return err
			}

			_handle(line)
		}

		return nil
	}else {
		return errors.New("类型错误.")
	}
}
