package models

type Bank struct {
	AuditModel
	Name        string `json:"name" gorm:"not null;size:50"`
	Description string `json:"description" gorm:"not null;size:255"`
	UserID      uint
	User        User `gorm:"foreignKey:UserID"`
}
