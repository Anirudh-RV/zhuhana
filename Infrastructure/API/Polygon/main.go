package main

import (
	"polygon/logger"
	"polygon/routes"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()
	log := logger.NewLogger()
	log.Info("Application started")

	routes.RegisterRoutes(router, log)

	router.Run("localhost:8080")
}
