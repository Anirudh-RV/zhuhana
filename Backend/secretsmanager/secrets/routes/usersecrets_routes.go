package routes

import (
	"database/sql"
	"secretsmanager/logger"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

func UserSecretsRoutesV1(r *gin.RouterGroup, log *logger.Logger, db *sql.DB, redis *redis.Client) {

}
