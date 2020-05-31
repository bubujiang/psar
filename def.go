package main

import (
	"flag"
	"fmt"
	"gopkg.in/ini.v1"
	"os"
)

type config struct {
	Ip string
	Port uint64
	PidFile string
}

func getCnf() *config {
	flag.String("c","conf.ini","配置文件路径")
	flag.Parse()

	cfg, err := ini.Load("my.ini")
	if err != nil {
		fmt.Printf("Fail to read file: %v", err)
		os.Exit(1)
	}

	cnf := &config{}
	cnf.Ip = cfg.Section("server").Key("ip").String()
	cnf.Port,_ = cfg.Section("server").Key("port").Uint64()
	cnf.PidFile = cfg.Section("server").Key("pid_file").String()

	return cnf
}
