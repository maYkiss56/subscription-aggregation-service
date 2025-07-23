package main

import (
	"log"

	"github.com/maYkiss56/subscription-aggregation-service/internal/app"
	"github.com/maYkiss56/subscription-aggregation-service/internal/config"
)

func main() {
	cfg := config.GetConfig()

	application := app.New(cfg)

	if err := application.Run(); err != nil {
		log.Fatalf("app failed: %v", err)
	}
}
