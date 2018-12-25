package main

import (
	"flag"

	"github.com/johncoleman83/cerebrum/pkg/api"

	"github.com/johncoleman83/cerebrum/pkg/utl/config"
)

func main() {
	cfg, err := config.LoadConfig()
	checkErr(err)

	checkErr(api.Start(cfg))
}

func checkErr(err error) {
	if err != nil {
		panic(err.Error())
	}
}
