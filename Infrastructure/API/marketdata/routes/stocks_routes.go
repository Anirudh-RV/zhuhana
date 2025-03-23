package routes

import (
	"marketdata/logger"
	tickersController "marketdata/stocks/controllers"
	tickersRepository "marketdata/stocks/repositories"
	tickersService "marketdata/stocks/services"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func StocksRoutesV1(r *gin.RouterGroup, log *logger.Logger) {
	stocks := r.Group("stocks/")
	{
		tickersRepo := tickersRepository.NewTickersRepository(log)
		log.Info("Tickers Repository created", zap.String("Execution Level", "Routes"))

		tickersService := tickersService.NewTickersService(tickersRepo, log)
		log.Info("Tickers Service created", zap.String("Execution Level", "Routes"))

		tickersController := tickersController.NewTickersController(tickersService, log)
		log.Info("Tickers Controller created", zap.String("Execution Level", "Routes"))

		stocks.GET("all-tickers/", tickersController.GetAllTickersV1)
	}
}
