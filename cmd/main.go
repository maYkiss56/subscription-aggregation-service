package main

import (
	"log"

	"github.com/maYkiss56/subscription-aggregation-service/internal/app"
	"github.com/maYkiss56/subscription-aggregation-service/internal/config"
)

// @title Subscription Aggregation Service API
// @version 1.0
// @description This is a service for managing user subscriptions
// @host localhost:8080
// @BasePath /api/subs
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
