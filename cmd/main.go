package main

import (
	"log"

	"github.com/maYkiss56/subscription-aggregation-service/internal/app"
	"github.com/maYkiss56/subscription-aggregation-service/internal/config"
)

func main() {
	cfg := config.GetConfig()

	application, err := app.New(cfg)
	if err != nil {
		log.Fatalf("failed to create app: %v", err)
	}

	if err := application.Run(); err != nil {
		log.Fatalf("app failed: %v", err)
	}
}
