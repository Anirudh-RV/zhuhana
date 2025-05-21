package email

import (
	"fmt"

	"go.uber.org/zap"
)

func (es *EmailService) SendOTPEmail(emailID, OTP, device, date, location, ipAddress string) error {
	go es.logger.Info("OTP for "+emailID+": "+OTP, zap.String("execution level", "SendOTPEmail"))
	if es.ENV == "dev" {
		return nil
	}

	subject := fmt.Sprintf(OTP_EMAIL_SUBJECT, OTP)
	return es.sendTemplatedEmail(
		emailID,
		subject,
		es.OTP_EMAIL_TEMPLATE_PATH,
		OTP, device, date, location, ipAddress,
	)
}
