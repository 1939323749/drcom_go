package main

import (
	"flag"
	"github.com/1939323749/drcom_go/conf"
	"github.com/1939323749/drcom_go/service"
	"log"
)

func main() {
	flag.Parse()
	if err := conf.Init(); err != nil {
		panic(err)
	}
	svr := service.New(conf.Conf)
	svr.Start()
	log.Printf("go-jlu-drcom-client [version: %s] start", conf.Conf.Version)
	c := make(chan struct{}, 0)
	<-c
}
