package main

import (
	"fmt"
	cli "github.com/jawher/mow.cli"
	"msar/boot"
	"msar/modules"
	"msar/modules/mem"
	"os"
)

func main() {

	serv := &server{}
	serv.run(getCnf())
	go watch()

	//app := cli.App("msar", "监控客户端")
	//app.Spec = "ACT"
	//src := app.StringArg("ACT", "restart", "start restart stop")
	//
	//app.Action = func() {
	//	switch (*src)[4:] {
	//	case "start":
	//		run()
	//		break
	//	case "stop":
	//		break
	//	default:
	//
	//	}
	//}
	//
	//err := app.Run(os.Args)
	//if err != err {}
}

func run(){
	//todo 优化
	p := &modules.Pack{}
	mem := modules.Module(&mem.Stats{})
	p.SetModule(&mem)
	p.SetType(mem.Type())
//	//go p.Run()
	go p.Run()

	for {
		x := <-boot.Data
		fmt.Println(&x)
		//squares <- x * x
	}
}
