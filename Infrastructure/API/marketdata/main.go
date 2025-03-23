package main

import (
	"marketdata/logger"
	"marketdata/routes"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func main() {
	router := gin.Default()
	log := logger.NewLogger()
	log.Info("Application started", zap.String("Execution Level", "Root"))

	routes.RegisterRoutes(router, log)

	router.Run("localhost:8080")
}
