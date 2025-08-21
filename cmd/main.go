package main

import (
	"context"
	"log"

	"github.com/nishant1479/Microservice-Backend/config"
)

func main() {
	log.Println("Initializing the infra services")
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	cfg, err := config.ParseConfig()
	if err != nil {
		log.Fatal(err)
	}

	appLogger := logger.NewApiLogger(cfg)
}