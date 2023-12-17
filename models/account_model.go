package models

type Account struct {
	AuditModel
	Name       string   `json:"name" gorm:"not null;size:100"`
	Balance    float64  `json:"balance" gorm:"type:decimal(15,2);not null;default:0.00"`
	CustomerID uint     `json:"customer_id" gorm:"not null"`
	Customer   Customer `json:"-" gorm:"foreignKey:CustomerID"`
}
