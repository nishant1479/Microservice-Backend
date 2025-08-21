package main

import (
	"context"
	"log"
)

func main() {
	log.Println("Initializing the infra services")
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	cfg, err := config.ParseConfig()
	if err != nil {
		log.Fatal(err)
	}
}