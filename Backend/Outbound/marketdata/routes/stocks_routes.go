package routes

import (
	"outbound/logger"
	polygonTickersController "outbound/marketdata/stocks/polygon/controllers"
	polygonTickersRepository "outbound/marketdata/stocks/polygon/repositories"
	polygonTickersService "outbound/marketdata/stocks/polygon/services"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

func StocksRoutesV1(r *gin.RouterGroup, log *logger.Logger, redis *redis.Client) {
	stocks := r.Group("stocks/")
	{
		polygon := stocks.Group("polygon/")
		{
			polygonTickersRepo := polygonTickersRepository.NewPolygonTickersRepository(log)
			go log.Info("Tickers Repository created", zap.String("Execution Level", "Routes"))

			polygonTickersService := polygonTickersService.NewPolygonTickersService(polygonTickersRepo, log)
			go log.Info("Tickers Service created", zap.String("Execution Level", "Routes"))

			polygonTickersController := polygonTickersController.NewPolygonTickersController(polygonTickersService, log)
			go log.Info("Tickers Controller created", zap.String("Execution Level", "Routes"))

			// TODO: Change full implementation for new Zhuana needs
			polygon.GET("ticker/", polygonTickersController.GetDailyTickerOHLCV_V1)
		}
	}
}
