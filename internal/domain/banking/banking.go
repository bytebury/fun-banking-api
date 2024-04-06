package banking

import (
	"funbanking/internal/domain"
	"funbanking/internal/domain/users"
	"funbanking/internal/infrastructure/persistence"
)

const (
	TransactionPending  = "pending"
	TransactionApproved = "approved"
	TransactionDeclined = "declined"
)

const EnablePremium = false

const (
	Free         int = iota // 0
	Premium                 // 1
	Family                  // 2
	Organization            // 3
)

type bankConfigLimits struct {
	Banks       int
	Employees   int
	Customers   int
	Accounts    int
	StoreFronts int
}

type bankSubscriptionTiers struct {
	Free         bankConfigLimits
	Premium      bankConfigLimits
	Family       bankConfigLimits
	Organization bankConfigLimits
}

type bankConfig struct {
	Limits bankSubscriptionTiers
}

var BankConfig = bankConfig{
	Limits: bankSubscriptionTiers{
		Free: bankConfigLimits{
			Banks:       2,
			Employees:   2,
			Customers:   25,
			Accounts:    2,
			StoreFronts: 0,
		},
		Premium: bankConfigLimits{
			Banks:       2,
			Employees:   2,
			Customers:   250,
			Accounts:    3,
			StoreFronts: 1,
		},
		Family: bankConfigLimits{
			Banks:       10,
			Employees:   10,
			Customers:   1_000,
			Accounts:    3,
			StoreFronts: 3,
		},
		Organization: bankConfigLimits{
			Banks:       100,
			Employees:   50,
			Customers:   10_000,
			Accounts:    3,
			StoreFronts: 5,
		},
	},
}

type Account struct {
	domain.AuditModel
	Name       string   `json:"name" gorm:"not null;size:100"`
	Balance    float64  `json:"balance" gorm:"type:decimal(50,2);not null;default:0.00"`
	CustomerID uint     `json:"customer_id" gorm:"not null"`
	Customer   Customer `json:"customer" gorm:"foreignKey:CustomerID;constraint:OnDelete:CASCADE;"`
	Type       string   `json:"type" gorm:"not null;default:checking"`
	IsPrimary  bool     `json:"is_primary" gorm:"not null;default:false"`
}

type AccountMonthlySummary struct {
	Month       string  `json:"month"`
	Deposits    float64 `json:"deposits"`
	Withdrawals float64 `json:"withdrawals"`
}

type Bank struct {
	domain.AuditModel
	Name        string     `json:"name" gorm:"not null;size:255"`
	Description string     `json:"description" gorm:"not null;size:255"`
	Slug        string     `json:"slug" gorm:"not null;size:255;uniqueIndex:idx_user_slug"`
	UserID      uint       `json:"user_id" gorm:"not null;uniqueIndex:idx_user_slug"`
	User        users.User `json:"owner" gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE;"`
}

type Customer struct {
	domain.AuditModel
	FirstName string    `json:"first_name" gorm:"not null;size:15"`
	LastName  string    `json:"last_name" gorm:"not null;size:20"`
	PIN       string    `json:"pin" gorm:"not null;default:000000;size:6;uniqueIndex:idx_pin_bank"`
	BankID    uint      `json:"bank_id" gorm:"not null;uniqueIndex:idx_pin_bank"`
	Bank      Bank      `json:"-" gorm:"foreignKey:BankID;constraint:OnDelete:CASCADE;"`
	Accounts  []Account `json:"accounts"`
}

type Employee struct {
	domain.AuditModel
	UserID uint       `json:"user_id" gorm:"not null;uniqueIndex:idx_user_bank"`
	User   users.User `json:"user" gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE;"`
	BankID uint       `json:"bank_id" gorm:"not null;uniqueIndex:idx_user_bank"`
	Bank   Bank       `json:"bank" gorm:"foreignKey:BankID;constraint:OnDelete:CASCADE;"`
}

type NewEmployeeRequest struct {
	BankID uint   `json:"bank_id"`
	Email  string `json:"email"`
}

type Transaction struct {
	domain.AuditModel
	Description       string     `json:"description" gorm:"not null;size:255"`
	CurrentBalance    float64    `json:"current_balance" gorm:"not null;type:decimal(50,2)"`
	Amount            float64    `json:"amount" gorm:"not null;type:decimal(50,2)"`
	Status            string     `json:"status" gorm:"not null;size:20;default:pending"`
	AccountID         uint       `json:"account_id" gorm:"not null"`
	Account           Account    `json:"account" gorm:"foreignKey:AccountID;constraint:OnDelete:CASCADE;"`
	UserID            *uint      `json:"user_id"`
	User              users.User `json:"user"`
	BankBuddySender   Customer   `json:"bank_buddy_sender"`
	BankBuddySenderID *uint      `json:"bank_buddy_sender_id"`
	Type              string     `json:"type" gorm:"not null;default:manual"`
}

type BankBuddyTransfer struct {
	FromAccountID uint    `json:"from_account_id"`
	ToAccountID   uint    `json:"to_account_id"`
	Amount        float64 `json:"amount"`
	Description   string  `json:"description"`
}

type TransferRequest struct {
	FromAccountID uint    `json:"from_account_id"`
	ToAccountID   uint    `json:"to_account_id"`
	Amount        float64 `json:"amount"`
}

func RunMigrations() {
	persistence.DB.AutoMigrate(&Employee{})
	persistence.DB.AutoMigrate(&Bank{})
	persistence.DB.AutoMigrate(&Customer{})
	persistence.DB.AutoMigrate(&Account{})
	persistence.DB.AutoMigrate(&Transaction{})
}
