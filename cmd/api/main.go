package main

import (
	"github.com/gin-gonic/gin"
	"github.com/omatheuscaetano/planus-api/internal/dependency_container"
	"github.com/omatheuscaetano/planus-api/internal/router"
	"github.com/omatheuscaetano/planus-api/pkg/env"
)

// @title Planus API
// @version 1.0
// @description This is the API documentation for the Planus application.
// @basePath /
func main() {
	env.Load()
	instances := dependency_container.InstantiateAll()

	gin := gin.Default()

	router.AllRoutes(gin, instances)

	gin.Run(":" + env.ApiPort())
}
