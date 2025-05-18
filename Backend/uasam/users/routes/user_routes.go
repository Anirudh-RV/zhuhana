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

	userController "uasam/users/user/controllers"
	userRepository "uasam/users/user/repositories"
	userService "uasam/users/user/services"
)

func UserRoutesV1(ctx *context.Context, r *gin.RouterGroup, log *logger.Logger, db *sql.DB, redis *redis.Client, emailService *email.EmailService, jwtService *commonutils.JWTService) {
	user := r.Group("user/")
	{
		userRepo := userRepository.NewUserRepository(db)
		go log.Info("user repository created", zap.String("execution level", "UserRoutesV1"))

		otpService := userService.NewOTPService(ctx, log, redis, emailService)
		go log.Info("otp service created", zap.String("execution level", "UserRoutesV1"))

		userService := userService.NewUserService(ctx, otpService, jwtService, userRepo, log, redis)
		go log.Info("user service created", zap.String("execution level", "UserRoutesV1"))

		signUp := user.Group("sign-up/")
		{
			signUpController := userController.NewSignUpController(userService, log)
			go log.Info("signup controller created", zap.String("execution level", "UserRoutesV1"))

			signUp.POST("init/", signUpController.SignUpInitHandler)
			signUp.POST("verify-otp/", signUpController.SignUpVerifyOTPHandler)
		}

		login := user.Group("login/")
		{
			loginController := userController.NewLoginController(userService, log)
			go log.Info("login controller created", zap.String("execution level", "UserRoutesV1"))

			login.POST("verify-password/", loginController.LoginVerifyPasswordHandler)
			login.POST("verify-otp/", loginController.LoginVerifyOTPHandler)
		}
	}
}
