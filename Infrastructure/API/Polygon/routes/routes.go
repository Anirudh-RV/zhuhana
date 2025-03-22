package routes

import (
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.Engine) {
	v1 := r.Group("/api/polygon/v1/")
	{
		StocksRoutesV1(v1)
	}
}
