package services

import (
	"golfer/models"
	"golfer/repositories"
	"strconv"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	userRepository repositories.UserRepository
	jwtService     JwtService
}

func NewUserService(userRepository repositories.UserRepository, jwtService JwtService) *UserService {
	return &UserService{
		userRepository,
		jwtService,
	}
}

func (service UserService) Create(request *models.UserRequest, out *models.User) error {
	passwordHash, err := service.HashPassword(request.Password)

	if err != nil {
		return err
	}

	service.normalize(request)

	out.Email = request.Email
	out.Username = request.Username
	out.FirstName = request.FirstName
	out.LastName = request.LastName
	out.Password = passwordHash

	return service.userRepository.Create(out)
}

func (service UserService) FindByID(userID string, user *models.User) error {
	return service.userRepository.FindByID(userID, user)
}

func (service UserService) FindByEmail(email string, user *models.User) error {
	return service.userRepository.FindByEmail(strings.ToLower(email), user)
}

func (service UserService) Update(userID string, request *models.UserRequest) (models.User, error) {
	var user models.User

	if err := service.userRepository.FindByID(userID, &user); err != nil {
		return user, err
	}

	service.normalize(request)

	if request.Username != "" {
		user.Username = request.Username
	}

	if request.Email != "" {
		user.Email = request.Email
	}

	if request.FirstName != "" {
		user.FirstName = request.FirstName
	}

	if request.LastName != "" {
		user.LastName = request.LastName
	}

	if err := service.userRepository.Update(&user); err != nil {
		return user, err
	}

	return user, nil
}

func (service UserService) Delete(userID string) error {
	return service.userRepository.Delete(userID)
}

func (service UserService) Login(request models.LoginRequest) (string, models.User, error) {
	var user models.User

	if err := service.FindByEmail(request.Email, &user); err != nil {
		return "", user, err
	}

	if !service.VerifyPassword(request.Password, user.Password) {
		return "", user, nil
	}

	token, err := service.jwtService.GenerateUserToken(strconv.Itoa(int(user.ID)))

	if err != nil {
		return "", user, err
	}
	return token, user, err
}

func (service UserService) HashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	if err != nil {
		return "", err
	}

	return string(hash), nil
}

func (service UserService) VerifyPassword(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

/**
 * Internal function used for services on the backend to update users.
 * If you are updating a user from an HTTP request, you should use Update().
 *
 * We do this to avoid passing the password "over-the-wire" so there is a
 * function to update the user on the backend, and one to deal with requests
 * to update a user.
 *
 * @see Update
 */
func (service UserService) updateWholeUser(user *models.User) error {
	return service.userRepository.Update(user)
}

func (service UserService) normalize(user *models.UserRequest) {
	user.Email = strings.ToLower(user.Email)
	user.Username = strings.ToLower(user.Username)
	user.FirstName = strings.ToLower(user.FirstName)
	user.LastName = strings.ToLower(user.LastName)
}
