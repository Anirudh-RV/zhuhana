package services

import (
	"context"
	"os"
	"strconv"
	"time"
	"uasam/email"
	"uasam/logger"
	"uasam/users/user/repositories"

	"github.com/pquerna/otp"
	"github.com/pquerna/otp/totp"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

type OTPService struct {
	ctx                                     *context.Context
	logger                                  *logger.Logger
	OTP_SKEW                                int
	OTP_ISSUER                              string
	OTP_SECRETS_KEY_SUFFIX                  string
	RESET_PASSWORD_KEY_SUFFIX               string
	OTP_DURATION_IN_SECONDS                 int
	RESET_PASSWORD_TOKEN_DURATION_IN_SECOND int
	redis                                   *redis.Client
	emailService                            *email.EmailService
	userRepository                          *repositories.UserRepository
}

func NewOTPService(ctx *context.Context, logger *logger.Logger, redis *redis.Client, emailService *email.EmailService, userRepository *repositories.UserRepository) *OTPService {
	OTP_SKEW, _ := strconv.Atoi(os.Getenv("OTP_SKEW"))
	OTP_SECRETS_KEY_SUFFIX := os.Getenv("OTP_SECRETS_KEY_SUFFIX")
	OTP_DURATION_IN_SECONDS, _ := strconv.Atoi(os.Getenv("OTP_DURATION_IN_SECONDS"))
	RESET_PASSWORD_KEY_SUFFIX := os.Getenv("RESET_PASSWORD_KEY_SUFFIX")
	RESET_PASSWORD_TOKEN_DURATION_IN_SECOND, _ := strconv.Atoi(os.Getenv("RESET_PASSWORD_TOKEN_DURATION_IN_SECOND"))

	return &OTPService{
		ctx:                                     ctx,
		logger:                                  logger,
		redis:                                   redis,
		emailService:                            emailService,
		userRepository:                          userRepository,
		OTP_SKEW:                                OTP_SKEW,
		OTP_ISSUER:                              os.Getenv("OTP_ISSUER"),
		OTP_SECRETS_KEY_SUFFIX:                  OTP_SECRETS_KEY_SUFFIX,
		OTP_DURATION_IN_SECONDS:                 OTP_DURATION_IN_SECONDS,
		RESET_PASSWORD_KEY_SUFFIX:               RESET_PASSWORD_KEY_SUFFIX,
		RESET_PASSWORD_TOKEN_DURATION_IN_SECOND: RESET_PASSWORD_TOKEN_DURATION_IN_SECOND,
	}
}

func (ots *OTPService) generateAndStoreSecretKey(emailID string) (string, error) {
	secret, err := totp.Generate(totp.GenerateOpts{
		Issuer:      ots.OTP_ISSUER,
		AccountName: emailID,
	})
	if err != nil {
		return "", err
	}
	err = ots.redis.Set(*ots.ctx, emailID+ots.OTP_SECRETS_KEY_SUFFIX, secret.Secret(), time.Duration(ots.OTP_DURATION_IN_SECONDS)*time.Second).Err()
	if err != nil {
		go ots.logger.Warning("could not store user secret in redis", zap.String("rxecution level", "generateAndStoreSecretKey"), zap.String("Error", err.Error()))
		return "", err
	}

	return secret.Secret(), nil
}

func (ots *OTPService) getStoredSecretKey(emailID string) (string, error) {
	secret, err := ots.redis.Get(*ots.ctx, emailID+ots.OTP_SECRETS_KEY_SUFFIX).Result()
	if err != nil {
		go ots.logger.Warning("Could not get User Secret from Redis", zap.String("execution level", "getStoredSecretKey"), zap.String("Error", err.Error()))
		return "", err
	}

	return secret, nil
}

func (ots *OTPService) generateOTP(secret string) (string, error) {
	code, err := totp.GenerateCode(secret, time.Now())
	if err != nil {
		return "", err
	}
	return code, nil
}

func (ots *OTPService) VerifyOTP(secret, userInput string) (bool, error) {

	opts := totp.ValidateOpts{
		Period:    30,                 // OTP changes every 30s
		Skew:      uint(ots.OTP_SKEW), // Accept ±10 intervals = 5 minutes total
		Digits:    otp.DigitsSix,      // 6-digit OTP
		Algorithm: otp.AlgorithmSHA1,
	}
	return totp.ValidateCustom(userInput, secret, time.Now(), opts)
}

func (ots *OTPService) SendOTP(emailID string) error {
	secret, err := ots.generateAndStoreSecretKey(emailID)
	if err != nil {
		return err
	}

	otp, err := ots.generateOTP(secret)
	if err != nil {
		return err
	}

	err = ots.emailService.SendOTPEmail(emailID, otp)
	if err != nil {
		return err
	}

	return nil
}
