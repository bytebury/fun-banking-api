package models

type Announcement struct {
	AuditModel
	Title       string `json:"title" gorm:"not null;size:100"`
	Description string `json:"description" gorm:"not null;type:text"`
	UserID      uint   `json:"user_id" gorm:"not null"`
	User        User   `json:"user" gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE;"`
}

type UpdateAnnouncementRequest struct {
	Title       string `json:"string"`
	Description string `json:"description"`
	UserID      uint   `json:"user_id"`
}
