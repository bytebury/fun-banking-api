package users

import (
	"errors"
	"strings"
)

type PasswordService interface {
	ForgotPassword(email string) error
	ResetPassword(email, password, passwordConfirmation string) error
}

type ForgotMailer interface {
	SendEmail(email string) error
}

type passwordService struct {
	forgotMailer       ForgotMailer
	passwordRepository PasswordRepository
	userRepository     UserRepository
}

func NewPasswordService(forgotMailer ForgotMailer) PasswordService {
	return passwordService{
		forgotMailer:       forgotMailer,
		passwordRepository: NewPasswordRepository(),
		userRepository:     NewUserRepository(),
	}
}

func (p passwordService) ForgotPassword(email string) error {
	if strings.TrimSpace(email) == "" {
		return errors.New("email is required")
	}

	var user User
	if err := p.userRepository.FindByUsernameOrEmail(email, &user); err != nil {
		return errors.New("user does not exist")
	}

	return p.forgotMailer.SendEmail(email)
}

func (p passwordService) ResetPassword(email, password, passwordConfirmation string) error {
	if password != passwordConfirmation {
		return errors.New("passwords do not match")
	}

	return p.passwordRepository.UpdatePassword(email, password)
}
