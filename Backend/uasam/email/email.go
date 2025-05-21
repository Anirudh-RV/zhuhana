package email

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"uasam/logger"

	"go.uber.org/zap"
)

type EmailService struct {
	ctx                                *context.Context
	logger                             *logger.Logger
	BREVO_API_KEY                      string
	BREVO_SEND_EMAIL_API_ENDPOINT      string
	SENDER_EMAIL_NAME                  string
	SENDER_EMAIL_ID                    string
	OTP_EMAIL_TEMPLATE_PATH            string
	RESET_PASSWORD_EMAIL_TEMPLATE_PATH string
	ENV                                string
}

func NewEmailService(ctx *context.Context, logger *logger.Logger) *EmailService {
	ENV := os.Getenv("ENV")
	BREVO_API_KEY := os.Getenv("BREVO_API_KEY")
	BREVO_SEND_EMAIL_API_ENDPOINT := os.Getenv("BREVO_SEND_EMAIL_API_ENDPOINT")
	SENDER_EMAIL_NAME := os.Getenv("SENDER_EMAIL_NAME")
	SENDER_EMAIL_ID := os.Getenv("SENDER_EMAIL_ID")

	return &EmailService{
		ctx:                                ctx,
		logger:                             logger,
		BREVO_API_KEY:                      BREVO_API_KEY,
		BREVO_SEND_EMAIL_API_ENDPOINT:      BREVO_SEND_EMAIL_API_ENDPOINT,
		SENDER_EMAIL_NAME:                  SENDER_EMAIL_NAME,
		SENDER_EMAIL_ID:                    SENDER_EMAIL_ID,
		OTP_EMAIL_TEMPLATE_PATH:            "email/templates/login_otp_template.html",
		RESET_PASSWORD_EMAIL_TEMPLATE_PATH: "email/templates/reset_password_template.html",
		ENV:                                ENV,
	}
}

func (es *EmailService) sendTemplatedEmail(
	emailID string,
	subject string,
	templateFile string,
	args ...interface{},
) error {
	go es.logger.Info("Sending email to "+emailID, zap.String("execution level", "sendTemplatedEmail"))

	cwd, _ := os.Getwd()
	templatePath := filepath.Join(cwd, templateFile)
	htmlBytes, err := os.ReadFile(templatePath)
	if err != nil {
		go es.logger.Error("error reading template file", zap.String("execution level", "sendTemplatedEmail"), zap.String("Error", err.Error()))
		return err
	}
	htmlContent := fmt.Sprintf(string(htmlBytes), args...)

	emailData := EmailRequest{
		Sender: SenderInfo{
			Name:  es.SENDER_EMAIL_NAME,
			Email: es.SENDER_EMAIL_ID,
		},
		To: []EmailAddress{
			{Email: emailID},
		},
		Subject:     subject,
		HTMLContent: htmlContent,
	}

	jsonData, err := json.Marshal(emailData)
	if err != nil {
		go es.logger.Error("error marshalling the email data", zap.String("execution level", "sendTemplatedEmail"), zap.String("Error", err.Error()))
		return err
	}

	req, err := http.NewRequest("POST", es.BREVO_SEND_EMAIL_API_ENDPOINT, bytes.NewBuffer(jsonData))
	if err != nil {
		go es.logger.Info("error creating request", zap.String("execution level", "sendTemplatedEmail"), zap.String("Error", err.Error()))
		return err
	}

	req.Header.Set("accept", "application/json")
	req.Header.Set("content-type", "application/json")
	req.Header.Set("api-key", es.BREVO_API_KEY)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		go es.logger.Info("error sending brevo api", zap.String("execution level", "sendTemplatedEmail"), zap.String("Error", err.Error()))
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		go es.logger.Info("issue with brevo endpoint response", zap.String("execution level", "sendTemplatedEmail"))
		return errors.New("error sending email with brevo")
	}

	return nil
}
