package announcements

import "funbanking/internal/domain/model"

type Announcement struct {
	model.AuditModel
	Title       string     `json:"title" gorm:"not null;size:100"`
	Description string     `json:"description" gorm:"not null;type:text"`
	UserID      uint       `json:"user_id" gorm:"not null"`
	User        model.User `json:"user" gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE;"`
}
