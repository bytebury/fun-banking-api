package services

import (
	"golfer/mailers"
	"golfer/models"
	"strings"
)

type PasswordService struct {
	userService  UserService
	jwt          JwtService
	forgotMailer mailers.ForgotPasswordMailer
}

func NewPasswordService(userService UserService, jwtService JwtService, forgotMailer mailers.ForgotPasswordMailer) *PasswordService {
	return &PasswordService{
		userService,
		jwtService,
		forgotMailer,
	}
}

func (service PasswordService) SendForgotPasswordEmail(recipient string) error {
	var user models.User

	recipient = strings.ToLower(recipient)

	if err := service.userService.FindByEmail(strings.ToLower(recipient), &user); err != nil {
		return err
	}

	token, err := service.jwt.GeneratePasswordResetToken(recipient)

	if err != nil {
		return err
	}

	return service.forgotMailer.SendEmail(recipient, token)
}

func (service PasswordService) UpdatePassword(email, password, confirmation string) error {
	var user models.User
	email = strings.ToLower(email)

	if err := service.userService.FindByEmail(email, &user); err != nil {
		return err
	}

	hashPassword, err := service.userService.HashPassword(password)

	if err != nil {
		return err
	}

	user.Password = hashPassword

	return service.userService.updateWholeUser(&user)
}
