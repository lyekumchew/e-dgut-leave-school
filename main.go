package main

import (
	"fmt"
	"github.com/lyekumchew/e-dgut-leave-school/config"
	"github.com/lyekumchew/e-dgut-leave-school/edgut"
	"os"
)

func main() {
	var conf config.Config
	err := conf.Get()
	if err != nil {
		fmt.Printf("conf error: #%v", err)
		os.Exit(1)
	}

	fmt.Println(conf)

	e := edgut.EDGUTClient{Config: conf}

	// login
	if err = e.Login(); err != nil {
		fmt.Printf("login error: #%v", err)
		os.Exit(1)
	}

	e.Do()
}
