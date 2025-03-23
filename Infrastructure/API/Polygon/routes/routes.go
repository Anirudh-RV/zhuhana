package routes

import (
	"polygon/logger"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.Engine, log *logger.Logger) {
	v1 := r.Group("/api/polygon/v1/")
	{
		StocksRoutesV1(v1, log)
	}
}
