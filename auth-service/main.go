package main

import (
	"auth-service/config"
	"auth-service/internal"
)

func main() {
	cfg, err := config.NewConfig()
	if err != nil {
		panic(err)
	}

	internal.Run(cfg)
}
