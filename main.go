package main

import (
	"flag"
	"fmt"
	"github.com/lyekumchew/e-dgut-leave-school/common"
	"github.com/lyekumchew/e-dgut-leave-school/config"
	"github.com/lyekumchew/e-dgut-leave-school/edgut"
	"github.com/robfig/cron/v3"
	"os"
)

var cronJob = flag.Bool("-cron", false, "是否启用 cronjob，默认为是，否则只运行一次")

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
	} else {
		do(conf)
	}
}
