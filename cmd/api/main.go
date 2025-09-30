package main

import (
	"context"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/omatheuscaetano/planus-api/internal/dependency_container"
	"github.com/omatheuscaetano/planus-api/internal/router"
	mongo "github.com/omatheuscaetano/planus-api/pkg/db"
	"github.com/omatheuscaetano/planus-api/pkg/env"
)

// @title Planus API
// @version 1.0
// @description This is the API documentation for the Planus application.
// @basePath /
func main() {
	env.Load()

	mongo := mongo.NewProductionMongoDb()

	err := mongo.Connect(context.Background())
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}

	instances := dependency_container.InstantiateAll(&dependency_container.InstancesConfig{
		Mongo: mongo,
	})

	gin := gin.Default()

	router.AllRoutes(gin, instances)

	gin.Run(":" + env.ApiPort())
}
