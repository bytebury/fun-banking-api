package repositories

import (
	"golfer/database"
	"golfer/models"
	"time"

	"gorm.io/gorm"
)

type HealthRepository struct {
	db *gorm.DB
}

func NewHealthRepository() *HealthRepository {
	return &HealthRepository{
		db: database.DB,
	}
}

func (repository HealthRepository) GetHealthCheck(health *models.Health) error {
	health.Name = "Fun Banking"
	health.Version = "1.0.0"
	health.Message = "Everything is up and running!"

	repository.db.Model(&models.User{}).Count(&health.NumberOfUsers)
	repository.db.Model(&models.Bank{}).Count(&health.NumberOfBanks)
	repository.db.Model(&models.Customer{}).Count(&health.NumberOfCustomers)
	repository.db.Model(&models.Transaction{}).Count(&health.NumberOfTransfers)

	return nil
}

func (repository HealthRepository) GetUserWeeklyInsights(insights *[]models.WeeklyInsights) error {
	twelveWeeksAgo := time.Now().AddDate(0, 0, -12*7)

	return repository.db.Model(&models.User{}).
		Select("EXTRACT(WEEK FROM created_at) as week, COUNT(*) as count").
		Where("created_at >= ?", twelveWeeksAgo).
		Group("week").
		Order("week").
		Scan(&insights).Error
}
