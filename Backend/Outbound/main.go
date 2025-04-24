package main

import (
	"outbound/logger"
	"outbound/routes"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func main() {
	router := gin.Default()
	log := logger.NewLogger()
	log.Info("Application started", zap.String("Execution Level", "Root"))

	routes.RegisterRoutes(router, log)

	router.Run(":8080")
}
