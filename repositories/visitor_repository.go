package repositories

import (
	"golfer/database"
	"golfer/models"
	"time"

	"gorm.io/gorm"
)

type VisitorRepository struct {
	db *gorm.DB
}

type VisitorByDay []struct {
	Date          time.Time `json:"date"`
	UserCount     int       `json:"user_count"`
	CustomerCount int       `json:"customer_count"`
	GuestCount    int       `json:"guest_count"`
	TotalCount    int       `json:"total_count"`
}

func NewVisitorRepository() *VisitorRepository {
	return &VisitorRepository{
		db: database.DB,
	}
}

func (r VisitorRepository) AddVisitor(visitor *models.Visitor) error {
	return r.db.Create(&visitor).Error
}

func (r VisitorRepository) GetVisitorsByDay(result *VisitorByDay) error {
	startDate := time.Now().AddDate(0, 0, -14)

	return r.db.Model(&models.Visitor{}).
		Select("DATE(created_at) as date, COUNT(DISTINCT CASE WHEN user_id IS NOT NULL THEN ip_address END) as user_count, COUNT(DISTINCT CASE WHEN customer_id IS NOT NULL THEN ip_address END) as customer_count, COUNT(DISTINCT CASE WHEN user_id IS NULL AND customer_id IS NULL THEN ip_address END) as guest_count, COUNT(DISTINCT ip_address) as total_count").
		Group("DATE(created_at)").
		Where("created_at >= ?", startDate).
		Order("date").
		Scan(&result).Error
}
