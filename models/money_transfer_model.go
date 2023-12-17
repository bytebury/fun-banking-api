package models

type MoneyTransfer struct {
	AuditModel
	Description    string  `json:"description" gorm:"not null;size:255"`
	CurrentBalance float64 `json:"current_balance" gorm:"not null;type:decimal(10,2)"`
	Amount         float64 `json:"amount" gorm:"not null;type:decimal(10,2)"`
	Status         string  `json:"status" gorm:"not null;size:20;default:pending"`
	AccountID      uint    `json:"account_id" gorm:"not null"`
	Account        Account `json:"-" gorm:"foreignKey:AccountID"`
	UserID         *uint   `json:"user_id"`
	User           User    `json:"updated_by"`
}
