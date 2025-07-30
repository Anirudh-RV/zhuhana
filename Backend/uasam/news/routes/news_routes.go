package routes

import (
	"context"
	"database/sql"
	"uasam/commonutils"
	"uasam/email"
	"uasam/logger"
	"uasam/middleware"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"

	newsControllers "uasam/news/controllers"
	newsServices "uasam/news/services"
)

func NewsRoutesV1(ctx *context.Context, r *gin.RouterGroup, log *logger.Logger, db *sql.DB, redis *redis.Client, emailService *email.EmailService, jwtService *commonutils.JWTService, authMiddleware gin.HandlerFunc, userAuthMiddleware gin.HandlerFunc) {
	news := r.Group("news/")
	{
		newsArticleService := newsServices.NewNewsArticleService(log)
		go log.Info("news article service created", zap.String("execution level", "NewsRoutesV1"))

		newsArticlesController := newsControllers.NewNewsArticlesController(log, newsArticleService)
		go log.Info("news article controller created", zap.String("execution level", "NewsRoutesV1"))

		news.GET("article/", middleware.RateLimiter(redis, log, middleware.RateLimiterConfig{
			Source:      "header",
			Param:       "USER_TOKEN",
			EnableParam: true,
			Limit:       50,
			Window:      300,
			EnableIP:    true,
			IPLimit:     100,
			IPWindow:    300,
			Endpoint:    "news/articles/",
		}), userAuthMiddleware, newsArticlesController.GetNewsArticle)

	}
}
