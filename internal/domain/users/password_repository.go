package users

import (
	"funbanking/internal/infrastructure/persistence"
	"strings"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type PasswordRepository interface {
	UpdatePassword(email, password string) error
}

type passwordRepository struct {
	db *gorm.DB
}

func NewPasswordRepository() PasswordRepository {
	return passwordRepository{db: persistence.DB}
}

func (r passwordRepository) UpdatePassword(email, password string) error {
	var user User

	if err := r.db.First(&user, "email = ?", strings.ToLower(email)).Error; err != nil {
		return err
	}

	passwordHash, err := HashString(password)

	if err != nil {
		return err
	}

	user.Password = passwordHash

	return r.db.Model(&user).Select("Password").Updates(&user).Error

}

func HashString(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	if err != nil {
		return "", err
	}

	return string(hash), nil
}
