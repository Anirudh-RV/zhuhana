package routes

import (
	logger "algonexus/logger"

	"github.com/gin-gonic/gin"

	"database/sql"

	"github.com/redis/go-redis/v9"
)

func RegisterTradeOrderManagerRoutesV1(
	r *gin.RouterGroup,
	log *logger.Logger,
	db *sql.DB,
	redis *redis.Client,
	auth gin.HandlerFunc,
) {
}
