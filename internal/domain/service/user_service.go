package service

import (
	"funbanking/internal/domain/model"
	"funbanking/internal/domain/repository"
	"funbanking/internal/infrastructure/auth"
	"funbanking/package/utils"
)

type UserService interface {
	FindByID(id string) (model.User, error)
	FindByUsernameOrEmail(usernameOrEmail string) (model.User, error)
	FindBanks(id string) ([]model.Bank, error)
	Update(id string, user *model.User) error
	Login(usernameOrEmail, password string) (string, model.User, error)
	Create(user *model.User) error
}

type userService struct {
	userRepository repository.UserRepository
}

func NewUserService(userRepository repository.UserRepository) UserService {
	return userService{userRepository}
}

func (s userService) FindByID(id string) (model.User, error) {
	var user model.User
	err := s.userRepository.FindByID(id, &user)
	return user, err
}

func (s userService) FindByUsernameOrEmail(usernameOrEmail string) (model.User, error) {
	var user model.User
	err := s.userRepository.FindByUsernameOrEmail(usernameOrEmail, &user)
	return user, err
}

func (s userService) FindBanks(id string) ([]model.Bank, error) {
	var banks []model.Bank
	err := s.userRepository.FindBanks(id, &banks)
	return utils.Listify(banks), err
}

func (s userService) Update(id string, user *model.User) error {
	return s.userRepository.Update(id, user)
}

func (s userService) Create(user *model.User) error {
	// TODO this will need to map a user to a new user request
	return s.userRepository.Create(user)
}

func (s userService) Login(usernameOrEmail, password string) (string, model.User, error) {
	authService := auth.NewUserAuth(s.userRepository)
	request := auth.LoginRequest{
		UsernameOrEmail: usernameOrEmail,
		Password:        password,
	}

	return authService.Login(request)
}
