package main

import (
	"OrderService/config"
	"OrderService/internal/app"
)

func main() {
	cfg, err := config.New()
	if err != nil {
		panic(err)
	}

	orderService := app.New(cfg)
	defer orderService.L.Sync()
	if err := orderService.Start(); err != nil {
		panic(err)
	}
}
