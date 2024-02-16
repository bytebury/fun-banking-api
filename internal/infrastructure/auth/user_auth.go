package auth

import (
	"errors"
	"funbanking/internal/domain/model"
	"funbanking/internal/domain/service"
	"strconv"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

type userAuth struct {
	userService service.UserService
	jwt         JWTService
}

func NewUserAuth(userService service.UserService) userAuth {
	return userAuth{
		userService: userService,
		jwt:         NewJWTService(),
	}
}

func (auth userAuth) Login(request struct {
	UsernameOrEmail string
	Password        string
}) (string, model.User, error) {
	request.UsernameOrEmail = strings.TrimSpace(strings.ToLower(request.UsernameOrEmail))

	user, err := auth.userService.FindByUsernameOrEmail(request.UsernameOrEmail)

	if err != nil {
		return "", user, err
	}

	if !verifyUserPassword(request.Password, user.Password) {
		return "", user, errors.New("invalid password")
	}

	token, err := auth.jwt.GenerateUserToken(strconv.Itoa(int(user.ID)))

	if err != nil {
		return "", user, err
	}

	return token, user, err
}

func HashString(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	if err != nil {
		return "", err
	}

	return string(hash), nil
}

func verifyUserPassword(password, hashedPassword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	return err == nil
}
