package main

import (
	"polygon/routes"

	"github.com/gin-gonic/gin"
)

func main() {
    router := gin.Default()

	routes.RegisterRoutes(router)

    router.Run("localhost:8080")
}
