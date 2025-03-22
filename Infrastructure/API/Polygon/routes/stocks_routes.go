package routes

import (
	"polygon/stocks"

	"github.com/gin-gonic/gin"
)

func StocksRoutesV1(r *gin.RouterGroup) {
	users := r.Group("stocks/")
	{
		users.GET("all-tickers/", stocks.GetAllTickersV1)
	}
}
