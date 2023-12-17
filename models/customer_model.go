package models

type Customer struct {
	AuditModel
	FirstName string `json:"first_name" gorm:"not null;size:15"`
	LastName  string `json:"last_name" gorm:"not null;size:20"`
	PIN       string `json:"pin" gorm:"not null;default:000000;size:6;uniqueIndex:idx_pin_bank"`
	BankID    uint   `json:"bank_id" gorm:"not null;uniqueIndex:idx_pin_bank"`
	Bank      Bank   `json:"bank" gorm:"foreignKey:BankID"`
}

type CustomerResponse struct {
	AuditModel
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	PIN       string `json:"pin"`
	BankID    uint   `json:"bank_id"`
}
