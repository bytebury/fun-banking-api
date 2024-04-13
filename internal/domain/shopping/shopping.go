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
	Item    []Item     `json:"items"`
}

type Item struct {
	domain.AuditModel
	Shop          Shop    `json:"shop" gorm:"not null;foreignKey:ShopID;constraint:OnDelete:CASCADE"`
	ShopID        uint    `json:"shop_id" gorm:"not null"`
	Name          string  `json:"name" gorm:"not null; size:100"`
	Description   string  `json:"description" gorm:"size:100"`
	NumberInStock int     `json:"number_in_stock" gorm:"not null; default: 0"`
	Price         float64 `json:"price" gorm:"not null;type:decimal(50,2)"`
}

type Purchase struct {
	domain.AuditModel
	CartID     string  `json:"cart_id" gorm:"not null;type:char(36)"`
	Item       Item    `json:"item" gorm:"foreignKey:ItemID;constraint:OnDelete:CASCADE"`
	ItemID     uint    `json:"item_id" gorm:"not null"`
	Price      float64 `json:"price" gorm:"not null;type:decimal(50,2)"`
	CartPrice  float64 `json:"cart_price" gorm:"not null;type:decimal(50,2)"`
	TotalTax   float64 `json:"total_tax" gorm:"not null;type:decimal(50,2)"`
	TotalPrice float64 `json:"total_price" gorm:"not null;type:decimal(50,2)"`
	TaxRate    float64 `json:"tax_rate" gorm:"not null; type:decimal(50,2)"`
}

func RunMigrations() {
	persistence.DB.AutoMigrate(&Shop{})
	persistence.DB.AutoMigrate(&Item{})
	persistence.DB.AutoMigrate(&Purchase{})
}
