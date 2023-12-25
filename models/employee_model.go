package models

type Employee struct {
	AuditModel
	UserID uint `json:"user_id" gorm:"not null"`
	User   User `json:"user" gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE;"`
	BankID uint `json:"bank_id" gorm:"not null"`
	Bank   Bank `json:"bank" gorm:"foreignKey:BankID;constraint:OnDelete:CASCADE;"`
}
