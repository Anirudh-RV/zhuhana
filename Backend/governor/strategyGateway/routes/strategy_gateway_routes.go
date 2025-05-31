package routes

import (
	"github.com/gin-gonic/gin"

	"governor/logger"
	orderManagerRoute "governor/strategyGateway/orderManager/routes"

	"database/sql"
	"github.com/redis/go-redis/v9"
)

func RegisterStrategyGatewayRoutesV1(
	r *gin.RouterGroup,
	log *logger.Logger,
	db *sql.DB,
	redis *redis.Client,
	auth gin.HandlerFunc,
) {
	strategyGateway := r.Group("/strategy")
	// orderManager.Use(auth)
	orderManagerRoute.RegisterTradeOrderManagerRoutesV1(strategyGateway, log, db, redis, auth)
}
