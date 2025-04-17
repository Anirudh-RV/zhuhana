package routes

import (
	"marketdata/logger"
	polygonTickersController "marketdata/stocks/polygon/controllers"
	polygonTickersRepository "marketdata/stocks/polygon/repositories"
	polygonTickersService "marketdata/stocks/polygon/services"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func StocksRoutesV1(r *gin.RouterGroup, log *logger.Logger) {
	stocks := r.Group("stocks/")
	{
		polygon := stocks.Group("polygon/")
		{
			polygonTickersRepo := polygonTickersRepository.NewPolygonTickersRepository(log)
			log.Info("Tickers Repository created", zap.String("Execution Level", "Routes"))

			polygonTickersService := polygonTickersService.NewPolygonTickersService(polygonTickersRepo, log)
			log.Info("Tickers Service created", zap.String("Execution Level", "Routes"))

			polygonTickersController := polygonTickersController.NewPolygonTickersController(polygonTickersService, log)
			log.Info("Tickers Controller created", zap.String("Execution Level", "Routes"))

			// TODO: Change full implementation for new Zhuana needs
			polygon.GET("all-tickers/", polygonTickersController.GetAllTickersV1)
		}
	}
}
