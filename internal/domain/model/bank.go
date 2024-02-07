package model

type Bank struct {
	AuditModel
	Name        string `json:"name" gorm:"not null;size:255"`
	Description string `json:"description" gorm:"not null;size:255"`
	Slug        string `json:"slug" gorm:"not null;size:255;uniqueIndex:idx_user_slug"`
	UserID      uint   `json:"user_id" gorm:"not null;uniqueIndex:idx_user_slug"`
	User        User   `json:"owner" gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE;"`
}
