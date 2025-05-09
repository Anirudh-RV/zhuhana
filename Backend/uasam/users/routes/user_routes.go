package routes

import (
	"context"
	"database/sql"
	"uasam/email"
	"uasam/logger"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"

	userController "uasam/users/user/controllers"
	userRepository "uasam/users/user/repositories"
	userService "uasam/users/user/services"
)

func UserRoutesV1(ctx *context.Context, r *gin.RouterGroup, log *logger.Logger, db *sql.DB, redis *redis.Client, emailService *email.EmailService) {
	user := r.Group("user/")
	{
		signUp := user.Group("sign-up/")
		{
			userRepo := userRepository.NewUserRepository(db)
			go log.Info("User Repository created", zap.String("Execution Level", "UserRoutesV1"))

			otpService := userService.NewOTPService(ctx, log, redis, emailService)
			go log.Info("OTP Service created", zap.String("Execution Level", "UserRoutesV1"))

			userService := userService.NewUserService(ctx, otpService, userRepo, log, redis)
			go log.Info("User Service created", zap.String("Execution Level", "UserRoutesV1"))

			signUpController := userController.NewSignUpController(userService, log)
			go log.Info("SignUp Controller created", zap.String("Execution Level", "UserRoutesV1"))

			signUp.POST("init/", signUpController.SignUpInitHandler)
			signUp.POST("verify-otp/", signUpController.VerifyOTPHandler)
		}
	}
}
