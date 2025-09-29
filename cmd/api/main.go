package main

import (
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

	mongodb := mongo.NewProductionMongoDb()
	instances := dependency_container.InstantiateAll(
		&dependency_container.InstancesConfig{
			Mongo: mongodb,
		},
	)

	gin := gin.Default()

	router.AllRoutes(gin, instances)

	gin.Run(":" + env.ApiPort())
}
