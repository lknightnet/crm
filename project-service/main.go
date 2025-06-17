package main

import (
	"project-service/config"
	"project-service/internal"
)

func main() {
	cfg, err := config.NewConfig()
	if err != nil {
		panic(err)
	}

	internal.Run(cfg)
}
