package users

import (
	"funbanking/internal/infrastructure/pagination"
	"funbanking/internal/infrastructure/persistence"

	"gorm.io/gorm"
)

type UserAuth interface {
	Login(request LoginRequest) (string, User, error)
}

type WelcomeMailer interface {
	SendEmail(recipient string, user User) error
}

type VerificationMailer interface {
	SendEmail(recipient string, user User) error
}

type LoginRequest struct {
	UsernameOrEmail string `json:"username_or_email"`
	Password        string `json:"password"`
}

type UserService interface {
	FindByID(id string) (User, error)
	FindByUsernameOrEmail(usernameOrEmail string) (User, error)
	FindAll(itemsPerPage, pageNumber int, params map[string]string) (pagination.PaginatedResponse[User], error)
	Update(id string, user *User) error
	Login(usernameOrEmail, password string) (string, User, error)
	Create(request *NewUserRequest) (User, error)
	AddVisitor(visitor *Visitor) error
	ChangeEmail(userID, email string) error
	Verify(email string) error
	ResendVerificationEmail(userID, email string) error
}

type userService struct {
	authService                   UserAuth
	userRepository                UserRepository
	welcomeMailer                 WelcomeMailer
	changeEmailVerificationMailer VerificationMailer
	accountVerificationMailer     VerificationMailer
}

func NewUserService(userRepository UserRepository, authService UserAuth, welcomeMailer WelcomeMailer, changeEmailVerificationMailer VerificationMailer, accountVerificationMailer VerificationMailer) UserService {
	return userService{
		userRepository:                userRepository,
		authService:                   authService,
		welcomeMailer:                 welcomeMailer,
		changeEmailVerificationMailer: changeEmailVerificationMailer,
		accountVerificationMailer:     accountVerificationMailer,
	}
}

func (s userService) FindByID(id string) (User, error) {
	var user User
	err := s.userRepository.FindByID(id, &user)
	return user, err
}

func (s userService) FindByUsernameOrEmail(usernameOrEmail string) (User, error) {
	var user User
	err := s.userRepository.FindByUsernameOrEmail(usernameOrEmail, &user)
	return user, err
}

func (s userService) FindAll(itemsPerPage, pageNumber int, params map[string]string) (pagination.PaginatedResponse[User], error) {
	return s.userRepository.FindAll(itemsPerPage, pageNumber, params)
}

func (s userService) Update(id string, user *User) error {
	return s.userRepository.Update(id, user)
}

func (s userService) Create(request *NewUserRequest) (User, error) {
	user := User{
		Username:  request.Username,
		Email:     request.Email,
		FirstName: request.FirstName,
		LastName:  request.LastName,
		Password:  request.Password,
		Role:      0,
		About:     "",
		Avatar:    "https://www.gravatar.com/avatar/2533c61da0bd2b79b63fd599cd045a31?default=https%3A%2F%2Fcloud.digitalocean.com%2Favatars%2Fdefault30.png&secure=true",
	}

	if err := s.userRepository.Create(&user); err != nil {
		return User{}, err
	}

	s.welcomeMailer.SendEmail(user.Email, user)

	return user, nil
}

func (s userService) Login(usernameOrEmail, password string) (string, User, error) {
	request := LoginRequest{
		UsernameOrEmail: usernameOrEmail,
		Password:        password,
	}

	return s.authService.Login(request)
}

func (s userService) AddVisitor(visitor *Visitor) error {
	return s.userRepository.AddVisitor(visitor)
}

func (s userService) ChangeEmail(userID, email string) error {
	return persistence.DB.Transaction(func(tx *gorm.DB) error {
		user, err := s.FindByID(userID)

		if err != nil {
			return err
		}

		user.Email = email

		if err := s.userRepository.UpdateEmail(&user); err != nil {
			return err
		}

		return s.changeEmailVerificationMailer.SendEmail(email, user)
	})
}

func (s userService) Verify(email string) error {
	user, err := s.FindByUsernameOrEmail(email)

	if err != nil {
		return err
	}

	user.Verified = true

	return s.userRepository.Verify(&user)
}

func (s userService) ResendVerificationEmail(userID, email string) error {
	if user, err := s.FindByID(userID); err != nil {
		return err
	} else {
		return s.accountVerificationMailer.SendEmail(email, user)
	}
}
