package routes

import (
	"database/sql"
	dockercontroller "forge/dockercontroller"
	"forge/logger"

	_ "forge/docs"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	PythonBuilderRoutesV1 "forge/userAlgorithmBuilder/pythonBuilder/routes"
)

func RegisterRoutes(r *gin.Engine, log *logger.Logger, db *sql.DB, redis *redis.Client, authMiddleware gin.HandlerFunc, dockerService *dockercontroller.DockerService) {
	v1 := r.Group("/v1/")
	{
		v1.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
		PythonBuilderRoutesV1.PythonBuilderRoutesV1(v1, log, db, redis, authMiddleware, dockerService)
	}
}
