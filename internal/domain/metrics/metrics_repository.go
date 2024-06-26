package metrics

import (
	"funbanking/internal/domain/banking"
	"funbanking/internal/domain/users"
	"funbanking/internal/infrastructure/persistence"
	"funbanking/package/utils"
	"time"

	"gorm.io/gorm"
)

type MetricRepository interface {
	GetApplicationInfo(appInfo *ApplicationInfo) error
	GetUsersInfo() ([]WeeklyInsights, error)
	GetVisitorsByDay() (VisitorByDay, error)
}

type metricRepository struct {
	db *gorm.DB
}

func NewMetricRepository() metricRepository {
	return metricRepository{db: persistence.DB}
}

func (r metricRepository) GetApplicationInfo(appInfo *ApplicationInfo) error {
	appInfo.Name = "Fun Banking"
	appInfo.Version = "1.0.0"
	appInfo.Message = "Everything is up and running!"

	r.db.Model(&users.User{}).Count(&appInfo.NumberOfUsers)
	r.db.Model(&banking.Bank{}).Count(&appInfo.NumberOfBanks)
	r.db.Model(&banking.Customer{}).Count(&appInfo.NumberOfCustomers)
	r.db.Model(&banking.Transaction{}).Count(&appInfo.NumberOfTransactions)

	return nil
}

func (r metricRepository) GetUsersInfo() ([]WeeklyInsights, error) {
	twelveWeeksAgo := time.Now().AddDate(0, 0, -12*7)
	var insights []WeeklyInsights

	err := r.db.Model(&users.User{}).
		Select("EXTRACT(WEEK FROM created_at) as week, COUNT(*) as count").
		Where("created_at >= ?", twelveWeeksAgo).
		Group("week").
		Order("week").
		Scan(&insights).Error

	return utils.Listify(insights), err
}

func (r metricRepository) GetVisitorsByDay() (VisitorByDay, error) {
	startDate := time.Now().AddDate(0, 0, -14)

	var result VisitorByDay

	err := r.db.Model(&users.Visitor{}).
		Select("DATE(created_at) as date, COUNT(DISTINCT CASE WHEN user_id IS NOT NULL THEN ip_address END) as user_count, COUNT(DISTINCT CASE WHEN customer_id IS NOT NULL THEN ip_address END) as customer_count, COUNT(DISTINCT ip_address) as total_count").
		Group("DATE(created_at)").
		Where("created_at >= ?", startDate).
		Order("date").
		Scan(&result).Error

	return result, err
}
