package models

type Customer struct {
	AuditModel
	FirstName string `json:"first_name" gorm:"not null;size:15"`
	LastName  string `json:"last_name" gorm:"not null;size:20"`
	PIN       string `json:"pin" gorm:"not null;size:6"`
	BankID    uint   `json:"bank_id"`
	Bank      Bank   `gorm:"foreignKey:BankID"`
}
