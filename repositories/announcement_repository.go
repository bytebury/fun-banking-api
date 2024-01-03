package repositories

import (
	"golfer/database"
	"golfer/models"

	"gorm.io/gorm"
)

type AnnouncementRepository struct {
	db *gorm.DB
}

func NewAnnouncementRepository() *AnnouncementRepository {
	return &AnnouncementRepository{
		db: database.DB,
	}
}

func (repository AnnouncementRepository) Create(announcement *models.Announcement) error {
	return repository.db.Create(&announcement).Error
}

func (repository AnnouncementRepository) Find(limit, page int) ([]models.Announcement, int64, error) {
	var announcements []models.Announcement
	var count int64
	repository.db.Find(&announcements).Count(&count)

	offset := (page - 1) * limit
	err := repository.db.Order("created_at DESC").Limit(limit).Offset(offset).Find(&announcements).Error

	return announcements, count, err
}

func (repository AnnouncementRepository) FindByID(announcementID string, announcement *models.Announcement) error {
	return repository.db.Preload("User").Find(&announcement, "id = ?", announcementID).Error
}

func (repository AnnouncementRepository) Update(announcement *models.Announcement) error {
	return repository.db.Save(&announcement).Error
}

func (repository AnnouncementRepository) Delete(announcementID string) error {
	return repository.db.Delete(&models.Announcement{}, "id = ?", announcementID).Error
}
