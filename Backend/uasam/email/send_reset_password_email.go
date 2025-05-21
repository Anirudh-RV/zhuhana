package email

import (
	"go.uber.org/zap"
)

func (es *EmailService) SendResetPasswordEmail(emailID, Token, device, date, location, ipAddress string) error {
	go es.logger.Info("Password Reset Token for "+emailID+": "+Token, zap.String("execution level", "SendResetPasswordEmail"))
	if es.ENV == "dev" {
		return nil
	}

	return es.sendTemplatedEmail(
		emailID,
		RESET_PASSWORD_SUBJECT,
		es.RESET_PASSWORD_EMAIL_TEMPLATE_PATH,
		Token, device, date, location, ipAddress,
	)
}
