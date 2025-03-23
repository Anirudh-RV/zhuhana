package routes

import (
	"polygon/logger"
	tickersController "polygon/stocks/controllers"
	tickersRepository "polygon/stocks/repositories"
	tickersService "polygon/stocks/services"

	"github.com/gin-gonic/gin"
)

func StocksRoutesV1(r *gin.RouterGroup, log *logger.Logger) {
	stocks := r.Group("stocks/")
	{
		tickersRepo := tickersRepository.NewTickersRepository(log)
		tickersService := tickersService.NewTickersService(tickersRepo, log)
		tickersController := tickersController.NewTickersController(tickersService, log)
		stocks.GET("all-tickers/", tickersController.GetAllTickersV1)
	}
}
