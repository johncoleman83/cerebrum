package main

import (
	"flag"

	"github.com/johncoleman83/cerebrum/pkg/api"

	"github.com/johncoleman83/cerebrum/pkg/utl/config"
)

func main() {
	errEnv := config.LoadEnvironment()
	checkErr(errEnv)
	
	cfgPath := flag.String("p", "./cmd/api/conf.local.yaml", "Path to config file")
	flag.Parse()

	cfg, err := config.Load(*cfgPath)
	checkErr(err)

	checkErr(api.Start(cfg))
}

func checkErr(err error) {
	if err != nil {
		panic(err.Error())
	}
}
