package mailing

import "fmt"

type ForgotPasswordMailer struct{}

func NewPasswordResetMailer() *ForgotPasswordMailer {
	return &ForgotPasswordMailer{}
}

func (mailer ForgotPasswordMailer) SendEmail(to, resetToken string) error {
	data := struct {
		AppName    string
		AppBaseURL string
		Token      string
	}{
		AppName:    "Fun Banking",
		AppBaseURL: "https://fun-banking.com",
		Token:      resetToken,
	}

	subject := fmt.Sprintf("%s - Reset Password Request", "Fun Banking")

	return sendEmail(to, subject, "password_reset", data)
}
