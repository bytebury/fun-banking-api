package models

type Bank struct {
	AuditModel
	Name        string `json:"name" gorm:"not null;size:50;`
	Description string `json:"description" gorm:"not null;size:255"`
	Slug        string `json:"slug" gorm:"not null;size:50;uniqueIndex:idx_user_slug"`
	UserID      uint   `json:"user_id" gorm:"not null;uniqueIndex:idx_user_slug"`
	User        User   `json:"owner" gorm:"foreignKey:UserID"`
}
