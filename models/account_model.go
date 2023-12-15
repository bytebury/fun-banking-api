package models

type Account struct {
	AuditModel,
	Name string `gorm:"not null;size:100"`
	Balance    float32  `gorm:"not null"`
	CustomerID uint     `json:"customer_id" gorm:"not null"`
	Customer   Customer `json:"customer" gorm:"foreignKey:CustomerID"`
}
