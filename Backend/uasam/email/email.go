package email

import (
	"context"
	"uasam/logger"

	"go.uber.org/zap"
)

type EmailService struct {
	ctx    *context.Context
	logger *logger.Logger
}

func NewEmailService(ctx *context.Context, logger *logger.Logger) *EmailService {
	return &EmailService{
		ctx:    ctx,
		logger: logger,
	}
}

func (es *EmailService) SendOTPEmail(emailID string, OTP string) error {
	// TODO: Implement Email OTP sender
	go es.logger.Info("OTP for "+emailID+": "+OTP, zap.String("execution level", "SendOTPEmail"))

	return nil
}

func (es *EmailService) SendResetPasswordEmail(emailID string, Token string) error {
	// TODO: Implement Email OTP sender
	go es.logger.Info("Password Reset Token for "+emailID+": "+Token, zap.String("execution level", "SendResetPasswordEmail"))

	return nil
}
