package repository

import (
	"funbanking/internal/domain/model"
	"funbanking/internal/infrastructure/persistence"
	"strings"

	"gorm.io/gorm"
)

type AnnouncementRepository interface {
	FindByID(id string, announcement *model.Announcement) error
	Create(announcement *model.Announcement) error
	Update(id string, announcement *model.Announcement) error
}

type announcementRepository struct {
	db *gorm.DB
}

func NewAnnouncementRepository() AnnouncementRepository {
	return announcementRepository{db: persistence.DB}
}

func (r announcementRepository) FindByID(id string, announcement *model.Announcement) error {
	return r.db.First(&announcement, "id = ?", id).Error
}

func (r announcementRepository) Create(annoncement *model.Announcement) error {
	return r.db.Create(&annoncement).Error
}

func (r announcementRepository) Update(id string, announcement *model.Announcement) error {
	var foundAnnouncement model.Announcement

	if err := r.FindByID(id, &foundAnnouncement); err != nil {
		return err
	}

	announcement.Title = strings.TrimSpace(announcement.Title)
	announcement.Description = strings.TrimSpace(announcement.Description)

	if announcement.Title == "" {
		announcement.Title = foundAnnouncement.Title
	}

	if announcement.Description == "" {
		announcement.Description = foundAnnouncement.Description
	}

	if announcement.UserID == 0 {
		announcement.UserID = foundAnnouncement.UserID
	}

	return r.db.Model(&foundAnnouncement).Select("Title", "Description", "UserID").Updates(&announcement).Error
}
