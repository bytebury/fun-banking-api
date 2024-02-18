package banking

import (
	"funbanking/internal/domain/model"
	"funbanking/internal/infrastructure/persistence"
)

const (
	TransactionPending  = "pending"
	TransactionApproved = "approved"
	TransactionDeclined = "declined"
)

type Account struct {
	model.AuditModel
	Name       string   `json:"name" gorm:"not null;size:100"`
	Balance    float64  `json:"balance" gorm:"type:decimal(50,2);not null;default:0.00"`
	CustomerID uint     `json:"customer_id" gorm:"not null"`
	Customer   Customer `json:"customer" gorm:"foreignKey:CustomerID;constraint:OnDelete:CASCADE;"`
}

type Bank struct {
	model.AuditModel
	Name        string     `json:"name" gorm:"not null;size:255"`
	Description string     `json:"description" gorm:"not null;size:255"`
	Slug        string     `json:"slug" gorm:"not null;size:255;uniqueIndex:idx_user_slug"`
	UserID      uint       `json:"user_id" gorm:"not null;uniqueIndex:idx_user_slug"`
	User        model.User `json:"owner" gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE;"`
}

type Customer struct {
	model.AuditModel
	FirstName string    `json:"first_name" gorm:"not null;size:15"`
	LastName  string    `json:"last_name" gorm:"not null;size:20"`
	PIN       string    `json:"pin" gorm:"not null;default:000000;size:6;uniqueIndex:idx_pin_bank"`
	BankID    uint      `json:"bank_id" gorm:"not null;uniqueIndex:idx_pin_bank"`
	Bank      Bank      `json:"-" gorm:"foreignKey:BankID;constraint:OnDelete:CASCADE;"`
	Accounts  []Account `json:"accounts"`
}

type Employee struct {
	model.AuditModel
	UserID uint       `json:"user_id" gorm:"not null;uniqueIndex:idx_user_bank"`
	User   model.User `json:"user" gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE;"`
	BankID uint       `json:"bank_id" gorm:"not null;uniqueIndex:idx_user_bank"`
	Bank   Bank       `json:"bank" gorm:"foreignKey:BankID;constraint:OnDelete:CASCADE;"`
}

type Transaction struct {
	model.AuditModel
	Description    string     `json:"description" gorm:"not null;size:255"`
	CurrentBalance float64    `json:"current_balance" gorm:"not null;type:decimal(50,2)"`
	Amount         float64    `json:"amount" gorm:"not null;type:decimal(50,2)"`
	Status         string     `json:"status" gorm:"not null;size:20;default:pending"`
	AccountID      uint       `json:"account_id" gorm:"not null"`
	Account        Account    `json:"account" gorm:"foreignKey:AccountID;constraint:OnDelete:CASCADE;"`
	UserID         *uint      `json:"user_id"`
	User           model.User `json:"updated_by"`
}

func RunMigrations() {
	persistence.DB.AutoMigrate(&Employee{})
	persistence.DB.AutoMigrate(&Bank{})
	persistence.DB.AutoMigrate(&Customer{})
	persistence.DB.AutoMigrate(&Account{})
	persistence.DB.AutoMigrate(&Transaction{})
}
