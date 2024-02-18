package metrics

import (
	"funbanking/internal/domain/banking"
	"funbanking/internal/domain/model"
	"funbanking/internal/infrastructure/persistence"

	"gorm.io/gorm"
)

type MetricRepository interface {
	GetApplicationInfo(appInfo *ApplicationInfo) error
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

	r.db.Model(&model.User{}).Count(&appInfo.NumberOfUsers)
	r.db.Model(&banking.Bank{}).Count(&appInfo.NumberOfBanks)
	r.db.Model(&banking.Customer{}).Count(&appInfo.NumberOfCustomers)
	r.db.Model(&banking.Transaction{}).Count(&appInfo.NumberOfTransactions)

	return nil
}
