package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/lyekumchew/e-dgut-leave-school/common"
	"github.com/lyekumchew/e-dgut-leave-school/config"
	"github.com/lyekumchew/e-dgut-leave-school/edgut"
	"github.com/robfig/cron/v3"
	"os"
	"os/signal"
	"syscall"
)

var cronJob = flag.Bool("cron", false, "enable cronjob")

func do(conf config.Config) {
	e := edgut.EDGUTClient{Config: conf}
	if err := e.Login(); err != nil {
		common.Logger(fmt.Sprintf("login error: #%v", err), 1)
		os.Exit(1)
	}
	err := e.Do()
	if err != nil {
		common.Logger(err.Error(), 1)
	}
}

func main() {
	flag.Parse()
	if *cronJob {
		common.Logger("开启每日自动离校请假", 2)
	} else {
		common.Logger("运行一次", 2)
	}

	// config
	var conf config.Config
	err := conf.Get()
	if err != nil {
		common.Logger(fmt.Sprintf("conf error: #%v", err), 1)
		os.Exit(1)
	}

	if *cronJob {
		c := cron.New()
		_, err = c.AddFunc("CRON_TZ=Asia/Shanghai 15 6 * * *", func() {
			do(conf)
		})
		if err != nil {
			panic("cronjob add function error")
		}
		c.Start()
		common.Logger("服务启动成功", 0)

		_, cancel := context.WithCancel(context.Background())
		defer cancel()
		ch := make(chan os.Signal, 2)
		signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
		select {
		case s := <-ch:
			cancel()
			common.Logger(fmt.Sprintf("\nreceived signal %s, exit.\n", s), 2)
		}
	} else {
		do(conf)
	}
}
