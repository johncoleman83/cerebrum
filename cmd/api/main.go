// Package main is the entry point to start the cerebrum server
package main

import (
	"github.com/johncoleman83/cerebrum/pkg/api"
	"github.com/johncoleman83/cerebrum/pkg/utl/config"
)

// main cerebrum server
func main() {
	cfg, err := config.LoadConfigFromFlags()
	if err != nil {
		panic(err.Error())
	}
	if cfg == nil {
		panic("unknown error loading yaml file")
	}

	if err = api.Start(cfg); err != nil {
		panic(err.Error())
	}
}
