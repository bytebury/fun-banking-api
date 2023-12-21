package models

type Customer struct {
	AuditModel
	FirstName string    `json:"first_name" gorm:"not null;size:15"`
	LastName  string    `json:"last_name" gorm:"not null;size:20"`
	PIN       string    `json:"pin" gorm:"not null;default:000000;size:6;uniqueIndex:idx_pin_bank"`
	BankID    uint      `json:"bank_id" gorm:"not null;uniqueIndex:idx_pin_bank"`
	Bank      Bank      `json:"-" gorm:"foreignKey:BankID;constraint:OnDelete:CASCADE;"`
	Accounts  []Account `json:"accounts"`
}

type CustomerSignInRequest struct {
	BankID string `json:"bank_id"`
	PIN    string `json:"pin"`
}
