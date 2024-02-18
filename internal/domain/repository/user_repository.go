package repository

import (
	"funbanking/internal/domain/model"
	"funbanking/internal/infrastructure/persistence"
	"strings"

	"gorm.io/gorm"
)

type UserRepository interface {
	GetCurrentUser(user *model.User) error
	FindByID(userID string, user *model.User) error
	FindByUsernameOrEmail(usernameOrEmail string, user *model.User) error
	Update(userID string, user *model.User) error
	Create(user *model.User) error
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository() UserRepository {
	return userRepository{db: persistence.DB}
}

func (r userRepository) GetCurrentUser(user *model.User) error {
	return r.db.Find(&user, "username = ?", "marcello").Error
}

func (r userRepository) FindByID(userID string, user *model.User) error {
	return r.db.Find(&user, "id = ?", userID).Error
}

func (r userRepository) FindByUsernameOrEmail(usernameOrEmail string, user *model.User) error {
	usernameOrEmail = strings.TrimSpace(strings.ToLower(usernameOrEmail))
	return r.db.Find(&user, "username = ? or email = ?", usernameOrEmail, usernameOrEmail).Error
}

func (r userRepository) Update(userID string, user *model.User) error {
	var foundUser model.User
	if err := r.FindByID(userID, &foundUser); err != nil {
		return err
	}

	if user.Username == "" {
		user.Username = foundUser.Username
	}

	if user.FirstName == "" {
		user.FirstName = foundUser.FirstName
	}

	if user.LastName == "" {
		user.LastName = foundUser.LastName
	}

	if user.Avatar == "" {
		user.Avatar = foundUser.Avatar
	}

	if user.About == "" {
		user.About = foundUser.About
	}

	r.normalize(user)

	return r.db.Model(&foundUser).Select("Username", "FirstName", "LastName", "Avatar", "About").Updates(&user).Error
}

func (r userRepository) Create(user *model.User) error {
	r.normalize(user)

	return r.db.Create(&user).Error
}

func (r userRepository) normalize(user *model.User) error {
	user.Username = strings.ToLower(user.Username)
	user.Email = strings.ToLower(user.Email)

	return nil
}
