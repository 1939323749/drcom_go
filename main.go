package main

import (
	"flag"
	"fmt"
	"github.com/1939323749/drcom_go/conf"
	"github.com/1939323749/drcom_go/service"
	"log"
)

func main() {
	c := make(chan struct{})
	flag.Parse()
	if err := conf.Init(); err != nil {
		panic(err)
	}
	svr := service.New(conf.Conf)
	svr.Start()

	defer func() {
		if arg := recover(); arg != nil {
			if err, ok := arg.(service.Error); ok {
				fmt.Println("Error:", err.Err)
				close(c)
			} else if str, ok := arg.(error); ok {
				fmt.Println("Message:", str)
			} else {
				fmt.Println("Panic recovered:", arg)
			}
		}
	}()

	log.Printf("go-jlu-drcom-client [version: %s] start", conf.Conf.Version)
	<-c
}
