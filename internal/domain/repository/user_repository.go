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
	Update(user *model.User) error
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

// TODO: Only update if it is present
func (r userRepository) Update(user *model.User) error {
	return r.db.Model(&user).Select("Username", "FirstName", "LastName", "Avatar", "About").Updates(&user).Error
}

func (r userRepository) Create(user *model.User) error {
	return r.db.Create(&user).Error
}
