package models

type PagingInfo struct {
	PageNumber uint `json:"page_number"`
	TotalItems uint `json:"total_items"`
}

type Transfer struct {
	AuditModel
	Description    string  `json:"description" gorm:"not null;size:255"`
	CurrentBalance float64 `json:"current_balance" gorm:"not null;type:decimal(50,2)"`
	Amount         float64 `json:"amount" gorm:"not null;type:decimal(50,2)"`
	Status         string  `json:"status" gorm:"not null;size:20;default:pending"`
	AccountID      uint    `json:"account_id" gorm:"not null"`
	Account        Account `json:"account" gorm:"foreignKey:AccountID;constraint:OnDelete:CASCADE;"`
	UserID         *uint   `json:"user_id"`
	User           User    `json:"updated_by"`
}

type PaginatedResponse[T any] struct {
	Items      []T        `json:"items"`
	PagingInfo PagingInfo `json:"paging_info"`
}
