package main

import (
	"allmygigs/config"
	"allmygigs/internal/db"
	"allmygigs/internal/infrastructure/routes"
	"log"

	"github.com/gin-gonic/gin"
)

func main() {
	cfg, err := config.LoadConfg()
	if err != nil {
		log.Fatalf("Error on config loading: %v", err)
	}

	dynamoClient, err := db.NewDynamoDBClient(cfg)
	if err != nil {
		log.Fatalf("Error on connecting database: %v", err)
	}

	if cfg.AppEnv == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	app := gin.Default()
	routes.SetupRoutes(app, cfg, dynamoClient)
	app.Run(":8080")
}
