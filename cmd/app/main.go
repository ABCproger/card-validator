package main

import (
	"log"

	"github.com/ABCproger/card-validator/config"
	"github.com/ABCproger/card-validator/internal/app"
)

func main() {
	cfg, err := config.New()
	if err != nil {
		log.Fatalf("config: %v", err)
	}
	app.Run(cfg)
}
