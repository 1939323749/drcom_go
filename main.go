package main

import (
	"flag"
	"github.com/1939323749/drcom_go/conf"
	"github.com/1939323749/drcom_go/service"
	"log"
	"os"
)

var infolog = log.New(os.Stdout, "[INFO]", log.LstdFlags)
var errlog = log.New(os.Stderr, "[ERROR]", log.LstdFlags)

func main() {
	flag.Parse()
	if err := conf.Init(); err != nil {
		panic(err)
	}
	svr := service.New(conf.Conf)

	if err := svr.Start(); err != nil {
		errlog.Printf("go-jlu-drcom-client [version: %s] start failed", conf.Conf.Version)
	} else {
		infolog.Printf("go-jlu-drcom-client [version: %s] started", conf.Conf.Version)
	}
	for {
		select {
		case _, ok := <-svr.LogoutCh:
			if !ok {
				svr.Restart <- 1
				//return
			}
		case restart, _ := <-svr.Restart:
			if restart > 3 {
				return
			}
			svr.Restart <- restart + 1
			err := svr.ReStart()
			if err != nil {
				return
			}
			infolog.Printf("go-jlu-drcom-client [version: %s] restarted, times:", conf.Conf.Version, restart)
		default:
		}
	}
}
