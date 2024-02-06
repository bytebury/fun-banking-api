package repository

import (
	"funbanking/internal/domain/model"
	"funbanking/internal/infrastructure/persistence"

	"gorm.io/gorm"
)

type User interface {
	GetCurrentUser(user *model.User) error
}

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository() UserRepository {
	return UserRepository{db: persistence.DB}
}

func (r UserRepository) GetCurrentUser(user *model.User) error {
	return r.db.Find(&user, "username = ?", "marcello").Error
}
