package models

type MoneyTransfer struct {
	AuditModel
	Description    string  `json:"description" gorm:"not null;size:255"`
	CurrentBalance float32 `json:"current_balance" gorm:"not null"`
	Amount         float32 `json:"amount" gorm:"not null"`
	Status         string  `json:"status" gorm:"not null;size:20;default:PENDING"`
	AccountID      uint    `json:"account_id" gorm:"not null"`
	Account        Account `json:"account" gorm:"foreignKey:CustomerID"`
}
