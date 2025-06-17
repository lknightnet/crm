package main

import (
	"jwt-service/config"
	"jwt-service/internal"
)

func main() {
	cfg, err := config.NewConfig()
	if err != nil {
		panic(err)
	}

	internal.Run(cfg)
}
