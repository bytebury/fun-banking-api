package shopping

import (
	"funbanking/internal/domain"
	"funbanking/internal/domain/users"
	"funbanking/internal/infrastructure/persistence"
)

type Shop struct {
	domain.AuditModel
	Name    string     `json:"name" gorm:"not null; size:100; uniqueIndex:idx_user_name"`
	TaxRate float64    `json:"tax_rate" gorm:"not null; default: 0.1"`
	UserID  uint       `json:"user_id" gorm:"not null; uniqueIndex:idx_user_name"`
	User    users.User `json:"user" gorm:"foreignKey:UserID; constraint:OnDelete:CASCADE"`
}

type Item struct {
	domain.AuditModel
	Name          string  `json:"name" gorm:"not null; size:100"`
	Description   string  `json:"description" gorm:"size:100"`
	NumberInStock int     `json:"number_in_stock" gorm:"not null; default: 0"`
	Price         float64 `json:"price" gorm:"not null;type:decimal(50,2)"`
}

type Purchase struct {
	domain.AuditModel
	GroupID    string  `json:"group_id" gorm:"not null;type:char(36)"`
	Item       Item    `json:"item" gorm:"foreignKey:ItemID;constraint:OnDelete:CASCADE"`
	ItemID     uint    `json:"item_id" gorm:"not null"`
	Price      float64 `json:"price" gorm:"not null;type:decimal(50,2)"`
	GroupPrice float64 `json:"group_price" gorm:"not null;type:decimal(50,2)"`
}

func RunMigrations() {
	persistence.DB.AutoMigrate(&Shop{})
	persistence.DB.AutoMigrate(&Item{})
	persistence.DB.AutoMigrate(&Purchase{})
}
