package auth

import (
	"errors"
	"funbanking/internal/domain/users"
	"strconv"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

type userAuth struct {
	userRepository users.UserRepository
	jwt            JWTService
}

func NewUserAuth(userRepository users.UserRepository) users.UserAuth {
	return userAuth{
		userRepository: userRepository,
		jwt:            NewJWTService(),
	}
}

func (auth userAuth) Login(request users.LoginRequest) (string, users.User, error) {
	request.UsernameOrEmail = strings.TrimSpace(strings.ToLower(request.UsernameOrEmail))

	var user users.User
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
