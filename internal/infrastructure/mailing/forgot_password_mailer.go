package mailing

import (
	"fmt"
	"funbanking/internal/infrastructure/auth"
)

type ForgotPasswordMailer struct {
	jwtService auth.JWTService
}

func NewPasswordResetMailer() *ForgotPasswordMailer {
	return &ForgotPasswordMailer{
		jwtService: auth.NewJWTService(),
	}
}

func (mailer ForgotPasswordMailer) SendEmail(recipient string) error {
	token, err := mailer.jwtService.GeneratePasswordResetToken(recipient)

	if err != nil {
		return err
	}

	data := struct {
		AppName    string
		AppBaseURL string
		Token      string
	}{
		AppName:    "Fun Banking",
		AppBaseURL: "https://fun-banking.com",
		Token:      token,
	}

	subject := fmt.Sprintf("%s - Reset Password Request", "Fun Banking")

	return sendEmail(recipient, subject, "forgot_password", data)
}
