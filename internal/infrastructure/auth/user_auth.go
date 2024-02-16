package auth

import (
	"errors"
	"funbanking/internal/domain/model"
	"funbanking/internal/domain/repository"
	"strconv"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

type LoginRequest struct {
	UsernameOrEmail string `json:"username_or_email"`
	Password        string `json:"password"`
}

type userAuth struct {
	userRepository repository.UserRepository
	jwt            JWTService
}

func NewUserAuth(userRepository repository.UserRepository) userAuth {
	return userAuth{
		userRepository: userRepository,
		jwt:            NewJWTService(),
	}
}

func (auth userAuth) Login(request LoginRequest) (string, model.User, error) {
	request.UsernameOrEmail = strings.TrimSpace(strings.ToLower(request.UsernameOrEmail))

	var user model.User
	if err := auth.userRepository.FindByUsernameOrEmail(request.UsernameOrEmail, &user); err != nil {
		return "", user, err
	}

	if !auth.verifyUserPassword(request.Password, user.Password) {
		return "", user, errors.New("invalid password")
	}

	token, err := auth.jwt.GenerateUserToken(strconv.Itoa(int(user.ID)))

	if err != nil {
		return "", user, err
	}

	return token, user, err
}

func (auth userAuth) verifyUserPassword(password, hashedPassword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	return err == nil
}
