package email

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"uasam/logger"

	"go.uber.org/zap"
)

type EmailService struct {
	ctx                           *context.Context
	logger                        *logger.Logger
	BREVO_API_KEY                 string
	BREVO_SEND_EMAIL_API_ENDPOINT string
	OTP_EMAIL_TEMPLATE            string
	OTP_EMAIL_SUBJECT             string
	SENDER_EMAIL_NAME             string
	SENDER_EMAIL_ID               string
}

type EmailBody struct {
}

func NewEmailService(ctx *context.Context, logger *logger.Logger) *EmailService {
	BREVO_API_KEY := os.Getenv("BREVO_API_KEY")
	BREVO_SEND_EMAIL_API_ENDPOINT := os.Getenv("BREVO_SEND_EMAIL_API_ENDPOINT")
	SENDER_EMAIL_NAME := os.Getenv("SENDER_EMAIL_NAME")
	SENDER_EMAIL_ID := os.Getenv("SENDER_EMAIL_ID")

	return &EmailService{
		ctx:                           ctx,
		logger:                        logger,
		BREVO_API_KEY:                 BREVO_API_KEY,
		BREVO_SEND_EMAIL_API_ENDPOINT: BREVO_SEND_EMAIL_API_ENDPOINT,
		OTP_EMAIL_TEMPLATE:            OTP_EMAIL_TEMPLATE,
		OTP_EMAIL_SUBJECT:             OTP_EMAIL_SUBJECT,
		SENDER_EMAIL_NAME:             SENDER_EMAIL_NAME,
		SENDER_EMAIL_ID:               SENDER_EMAIL_ID,
	}
}

func (es *EmailService) SendOTPEmail(emailID, OTP, device, date, location, ipAddress string) error {
	go es.logger.Info("OTP for "+emailID+": "+OTP, zap.String("execution level", "SendOTPEmail"))
	emailData := EmailRequest{
		Sender: SenderInfo{
			Name:  es.SENDER_EMAIL_NAME,
			Email: es.SENDER_EMAIL_ID,
		},
		To: []EmailAddress{
			{Email: emailID},
		},
		Subject:     fmt.Sprintf(es.OTP_EMAIL_SUBJECT, OTP),
		HTMLContent: fmt.Sprintf(es.OTP_EMAIL_TEMPLATE, OTP, device, date, location, ipAddress),
	}

	jsonData, err := json.Marshal(emailData)
	if err != nil {
		go es.logger.Error("error marshalling the email data", zap.String("execution level", "SendOTPEmail"), zap.String("Error", err.Error()))
		return err
	}

	// Create HTTP request
	req, err := http.NewRequest("POST", es.BREVO_SEND_EMAIL_API_ENDPOINT, bytes.NewBuffer(jsonData))
	if err != nil {
		go es.logger.Info("error creating request", zap.String("execution level", "SEND_OTP_EMAIL"), zap.String("Error", err.Error()))
		return err
	}

	// Add headers
	req.Header.Set("accept", "application/json")
	req.Header.Set("content-type", "application/json")
	req.Header.Set("api-key", es.BREVO_API_KEY) // Make sure to set this env var

	// Send request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		go es.logger.Info("error sending brevo api", zap.String("execution level", "SendOTPEmail"), zap.String("Error", err.Error()))
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		go es.logger.Info("issue with brevo endpoint response", zap.String("execution level", "SendOTPEmail"))
		return errors.New("error sending email with brevo")
	}

	return nil
}

func (es *EmailService) SendResetPasswordEmail(emailID string, Token string) error {
	// TODO: Implement Email OTP sender
	go es.logger.Info("Password Reset Token for "+emailID+": "+Token, zap.String("execution level", "SendResetPasswordEmail"))

	return nil
}
