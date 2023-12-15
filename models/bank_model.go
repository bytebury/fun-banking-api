package models

type Bank struct {
	AuditModel
	Name        string `json:"name" gorm:"not null;size:50;uniqueIndex:idx_user_name_slug"`
	Description string `json:"description" gorm:"not null;size:255"`
	Slug        string `json:"slug" gorm:"not null;size:50;uniqueIndex:idx_user_name_slug"`
	UserID      uint   `gorm:"not null;uniqueIndex:idx_user_name_slug"`
	User        User   `gorm:"foreignKey:UserID"`
}
