package repository

import (
	"funbanking/internal/domain/model"
	"funbanking/internal/infrastructure/persistence"

	"gorm.io/gorm"
)

type UserRepository interface {
	GetCurrentUser(user *model.User) error
	FindByID(userID string, user *model.User) error
	FindByUsernameOrEmail(usernameOrEmail string, user *model.User) error
	FindBanks(id string, banks *[]model.Bank) error
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
	return r.db.Find(&user, "username = ? or email = ?", usernameOrEmail).Error
}

func (r userRepository) FindBanks(id string, banks *[]model.Bank) error {
	return r.db.Find(&banks, "user_id = ?", id).Error
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

	return r.db.Model(&foundUser).Select("Username", "FirstName", "LastName", "Avatar", "About").Updates(&user).Error
}

func (r userRepository) Create(user *model.User) error {
	return r.db.Create(&user).Error
}
