package routes

import (
	"database/sql"
	"uasam/commonutils"
	"uasam/logger"
	microServiceController "uasam/microservices/microservice/controllers"
	microServiceService "uasam/microservices/microservice/services"
	"uasam/middleware"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

func MicroServiceRoutesV1(r *gin.RouterGroup, log *logger.Logger, db *sql.DB, redis *redis.Client, jwtService *commonutils.JWTService, authMiddleware gin.HandlerFunc) {
	microservicesRoute := r.Group("microservice/")
	{
		microServiceServiceObj := microServiceService.NewMicroServiceService(log, jwtService)
		go log.Info("microservice service created", zap.String("execution level", "MicroServiceRoutesV1"))

		microServiceLoginController := microServiceController.NewMicroServiceLoginController(microServiceServiceObj, log)
		go log.Info("microservice login controller created", zap.String("execution level", "MicroServiceRoutesV1"))

		microServiceAuthenticateController := microServiceController.NewMicroServiceAuthenticateController(microServiceServiceObj, log)
		go log.Info("microservice login controller created", zap.String("execution level", "MicroServiceRoutesV1"))

		microservicesRoute.POST("login/", microServiceLoginController.MicroServiceLoginHandler)
		microservicesRoute.POST("authenticate/", microServiceAuthenticateController.MicroServiceAuthenticateHandler)

		user := microservicesRoute.Group("user/")
		{
			microServiceUserLoginController := microServiceController.NewMicroServiceUserLoginController(microServiceServiceObj, log)
			go log.Info("microservice user login controller created", zap.String("execution level", "MicroServiceRoutesV1"))

			microServiceUserAuthenticateController := microServiceController.NewMicroServiceUserAuthenticateController(microServiceServiceObj, log)
			go log.Info("microservice user login controller created", zap.String("execution level", "MicroServiceRoutesV1"))

			user.POST("login/", middleware.RateLimiter(redis, log, middleware.RateLimiterConfig{
				Source:      "body",
				Param:       "userID",
				EnableParam: true,
				Limit:       3,
				Window:      300,
				EnableIP:    false,
				Endpoint:    "microservice/user/login/",
			}), authMiddleware, microServiceUserLoginController.MicroServiceUserLoginHandler)

			user.POST("authenticate/", middleware.RateLimiter(redis, log, middleware.RateLimiterConfig{
				Source:      "header",
				Param:       "USER_SERVICE_TOKEN",
				EnableParam: true,
				Limit:       300,
				Window:      300,
				EnableIP:    false,
				Endpoint:    "microservice/user/authenticate/",
			}), authMiddleware, microServiceUserAuthenticateController.MicroServiceUserAuthenticateHandler)
		}
	}
}
