package routes

import (
	logger "algonexus/logger"
	"algonexus/middleware"
	"algonexus/ordermanager/backtestengine/controllers"
	"algonexus/ordermanager/backtestengine/repositories"
	"algonexus/ordermanager/backtestengine/services"

	"go.uber.org/zap"

	"github.com/ClickHouse/clickhouse-go/v2"
	"github.com/gin-gonic/gin"

	"database/sql"

	"github.com/redis/go-redis/v9"
)

func RegisterBacktestRoutesV1(
	r *gin.RouterGroup,
	log *logger.Logger,
	db *sql.DB,
	clickHouse *clickhouse.Conn,
	redis *redis.Client,
	auth gin.HandlerFunc,
) {

	// Manager service init
	backtest := r.Group("backtest/")
	{
		backtestRepo := repositories.NewBacktestRepository(clickHouse)
		go log.Info("backtest repository created", zap.String("execution level", "RegisterBacktestRoutesV1"))

		backtestService := services.NewBacktestService(log, clickHouse, backtestRepo)
		go log.Info("backtest service created", zap.String("execution level", "UserRoutesV1"))

		backtestController := controllers.NewBacktestController(backtestService, log)
		go log.Info("backtest controller created", zap.String("execution level", "UserRoutesV1"))

		ohlc := backtest.Group("ohlc/")
		{
			ohlc.GET("range/", middleware.RateLimiter(redis, log, middleware.RateLimiterConfig{
				Source:      "header",
				Param:       "USER_ALGORITHM_TOKEN",
				EnableParam: true,
				Limit:       3,
				Window:      300,
				EnableIP:    true,
				IPLimit:     15,
				IPWindow:    300,
				Endpoint:    "/v1/backtest/ohlc/range/",
			}), backtestController.GetOHLCDataWithRange)
		}
	}
}
