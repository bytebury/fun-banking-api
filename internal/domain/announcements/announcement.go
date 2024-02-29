package announcements

import (
	"funbanking/internal/domain"
	"funbanking/internal/domain/users"
)

type Announcement struct {
	domain.AuditModel
	Title       string     `json:"title" gorm:"not null;size:100"`
	Description string     `json:"description" gorm:"not null;type:text"`
	UserID      uint       `json:"user_id" gorm:"not null"`
	User        users.User `json:"user" gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE;"`
}
