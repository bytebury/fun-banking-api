package repositories

import (
	"golfer/database"
	"golfer/models"

	"gorm.io/gorm"
)

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository() *UserRepository {
	return &UserRepository{
		db: database.DB,
	}
}

func (repository UserRepository) Create(user *models.User) error {
	return repository.db.Create(&user).Error
}

func (repository UserRepository) FindByID(userID string, user *models.User) error {
	return repository.db.First(&user, "id = ?", userID).Error
}

func (repository UserRepository) FindByEmail(email string, user *models.User) error {
	return repository.db.First(&user, "email = ?", email).Error
}

func (repository UserRepository) FindByUsername(username string, user *models.User) error {
	return repository.db.First(&user, "username = ?", username).Error
}

func (repository UserRepository) Update(user *models.User) error {
	return repository.db.Save(&user).Error
}

func (repository UserRepository) Delete(userId string) error {
	return repository.db.Delete(&models.User{}, "id = ?", userId).Error
}
