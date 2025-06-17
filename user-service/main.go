package main

import (
	"user-service/config"
	"user-service/internal"
)

func main() {
	cfg, err := config.NewConfig()
	if err != nil {
		panic(err)
	}

	internal.Run(cfg)
}
