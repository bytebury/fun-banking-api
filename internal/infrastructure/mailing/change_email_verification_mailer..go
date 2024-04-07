package mailing

import (
	"fmt"
	"funbanking/internal/domain/users"
	"funbanking/internal/infrastructure/auth"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

type ChangeEmailVerificationMailer struct {
	jwtService auth.JWTService
}

func NewChangeEmailVerificationMailer() *ChangeEmailVerificationMailer {
	return &ChangeEmailVerificationMailer{
		jwtService: auth.NewJWTService(),
	}
}

func (mailer ChangeEmailVerificationMailer) SendEmail(recipient string, user users.User) error {
	token, err := mailer.jwtService.GenerateVerificationToken(recipient)

	if err != nil {
		return err
	}

	user.FirstName = cases.Title(language.AmericanEnglish).String(user.FirstName)

	data := struct {
		User       users.User
		AppName    string
		AppBaseURL string
		Token      string
	}{
		User:       user,
		AppName:    "Fun Banking",
		AppBaseURL: "https://fun-banking.com",
		Token:      token,
	}

	subject := fmt.Sprintf("%s - E-mail Change Verification", "Fun Banking")

	return sendEmail(recipient, subject, "change_email_verification", data)
}
