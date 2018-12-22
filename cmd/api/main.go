package main

import (
	"flag"
	"os"
	"fmt"

	"github.com/johncoleman83/cerebrum/pkg/api"

	"github.com/johncoleman83/cerebrum/pkg/utl/config"
)

func main() {
	envs := map[string]bool{
    "dev": true,
    "stg": true,
		"prod": true,
	}
	// set ENVIRONMENT_NAME to dev if unset
	if _, err := envs["ENVIRONMENT_NAME"]; err {
    os.Setenv("ENVIRONMENT_NAME", "dev")
	}
	fmt.Println(os.Getenv("ENVIRONMENT_NAME"))
	
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
