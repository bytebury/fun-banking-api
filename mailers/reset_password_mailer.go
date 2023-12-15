package mailers

import (
	"fmt"
	"golfer/config"
)

type ForgotPasswordMailer struct{}

func NewPasswordResetMailer() *ForgotPasswordMailer {
	return &ForgotPasswordMailer{}
}

func (mailer ForgotPasswordMailer) SendEmail(to, resetToken string) error {
	const templateName = "password_reset"

	data := struct {
		AppName    string
		AppBaseURL string
		Token      string
	}{
		AppName:    config.AppName,
		AppBaseURL: config.AppBaseURL,
		Token:      resetToken,
	}

	subject := fmt.Sprintf("%s - Reset Password Request", config.AppName)

	return sendEmail(to, subject, templateName, data)
}
