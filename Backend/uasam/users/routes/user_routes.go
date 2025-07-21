package routes

import (
	"context"
	"database/sql"
	"uasam/commonutils"
	"uasam/email"
	"uasam/logger"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"

	middleware "uasam/middleware"
	userController "uasam/users/user/controllers"
	userRepository "uasam/users/user/repositories"
	userService "uasam/users/user/services"
)

func UserRoutesV1(ctx *context.Context, r *gin.RouterGroup, log *logger.Logger, db *sql.DB, redis *redis.Client, emailService *email.EmailService, jwtService *commonutils.JWTService, authMiddleware gin.HandlerFunc) {
	user := r.Group("user/")
	{
		userRepo := userRepository.NewUserRepository(db)
		go log.Info("user repository created", zap.String("execution level", "UserRoutesV1"))

		otpService := userService.NewOTPService(ctx, log, redis, emailService, userRepo)
		go log.Info("otp service created", zap.String("execution level", "UserRoutesV1"))

		userService := userService.NewUserService(ctx, otpService, jwtService, userRepo, log, redis)
		go log.Info("user service created", zap.String("execution level", "UserRoutesV1"))

		userAuthenticateController := userController.NewUserAuthenticateController(log, userService)
		go log.Info("user service created", zap.String("execution level", "UserRoutesV1"))

		signUp := user.Group("sign-up/")
		{
			signUpController := userController.NewSignUpController(userService, log)
			go log.Info("signup controller created", zap.String("execution level", "UserRoutesV1"))

			signUp.POST("init/", middleware.RateLimiter(redis, log, middleware.RateLimiterConfig{
				Source:      "body",
				Param:       "emailId",
				EnableParam: true,
				Limit:       3,
				Window:      300,
				EnableIP:    true,
				IPLimit:     15,
				IPWindow:    300,
				Endpoint:    "user/sign-up/init",
			}), signUpController.SignUpInitHandler)
			signUp.POST("verify-otp/", middleware.RateLimiter(redis, log, middleware.RateLimiterConfig{
				Source:      "body",
				Param:       "emailId",
				EnableParam: true,
				Limit:       3,
				Window:      300,
				EnableIP:    true,
				IPLimit:     15,
				IPWindow:    300,
				Endpoint:    "user/sign-up/verify-otp",
			}), signUpController.SignUpVerifyOTPHandler)
		}

		login := user.Group("login/")
		{
			loginController := userController.NewLoginController(userService, log)
			go log.Info("login controller created", zap.String("execution level", "UserRoutesV1"))

			login.POST("verify-password/", middleware.RateLimiter(redis, log, middleware.RateLimiterConfig{
				Source:      "body",
				Param:       "emailId",
				EnableParam: true,
				Limit:       10,
				Window:      300,
				EnableIP:    true,
				IPLimit:     15,
				IPWindow:    300,
				Endpoint:    "user/login/verify-password",
			}), loginController.LoginVerifyPasswordHandler)
			login.POST("verify-otp/", middleware.RateLimiter(redis, log, middleware.RateLimiterConfig{
				Source:      "body",
				Param:       "emailId",
				EnableParam: true,
				Limit:       3,
				Window:      300,
				EnableIP:    true,
				IPLimit:     15,
				IPWindow:    300,
				Endpoint:    "user/login/verify-otp",
			}), loginController.LoginVerifyOTPHandler)
		}

		resetPassword := user.Group("reset-password/")
		{
			resetPasswordController := userController.NewResetPasswordController(userService, log)
			go log.Info("reset password controller created", zap.String("execution level", "UserRoutesV1"))

			resetPassword.POST("init/", middleware.RateLimiter(redis, log, middleware.RateLimiterConfig{
				Source:      "body",
				Param:       "emailId",
				EnableParam: true,
				Limit:       3,
				Window:      300,
				EnableIP:    true,
				IPLimit:     15,
				IPWindow:    300,
				Endpoint:    "user/reset-password/init",
			}), resetPasswordController.ResetPasswordInitHandler)
			resetPassword.POST("reset/", middleware.RateLimiter(redis, log, middleware.RateLimiterConfig{
				Source:      "body",
				Param:       "emailId",
				EnableParam: true,
				Limit:       3,
				Window:      300,
				EnableIP:    true,
				IPLimit:     15,
				IPWindow:    300,
				Endpoint:    "user/reset-password/reset",
			}), resetPasswordController.ResetPasswordHandler)
		}

		user.POST("authenticate/", middleware.RateLimiter(redis, log, middleware.RateLimiterConfig{
			Source:      "header",
			Param:       "USER_TOKEN",
			EnableParam: true,
			Limit:       300,
			Window:      300,
			EnableIP:    true,
			IPLimit:     300,
			IPWindow:    300,
			Endpoint:    "/v1/user/authenticate/",
		}), userAuthenticateController.UserAuthenticateHandler)
	}
}
