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
	go es.logger.Info("OTP for "+emailID+": "+OTP, zap.String("Execution Level", "SendOTPEmail"))

	return nil
}
