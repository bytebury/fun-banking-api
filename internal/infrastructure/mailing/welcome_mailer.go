package mailing

import (
	"fmt"
	"funbanking/internal/domain/users"
	"funbanking/internal/infrastructure/auth"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

type WelcomeMailer struct {
	jwtService auth.JWTService
}

func NewWelcomeMailer() *WelcomeMailer {
	return &WelcomeMailer{
		jwtService: auth.NewJWTService(),
	}
}

func (mailer WelcomeMailer) SendEmail(recipient string, user users.User) error {
	user.FirstName = cases.Title(language.AmericanEnglish).String(user.FirstName)

	data := struct {
		User users.User
	}{
		User: user,
	}

	subject := fmt.Sprintf("%s - Welcome to Fun Banking 🎉", "Fun Banking")

	return sendEmail(recipient, subject, "welcome", data)
}
