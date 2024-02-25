package announcements

import (
	"funbanking/internal/infrastructure/pagination"
	"funbanking/internal/infrastructure/persistence"
	"strings"

	"gorm.io/gorm"
)

type AnnouncementRepository interface {
	FindByID(id string, announcement *Announcement) error
	FindAll(itemsPerPage, pageNumber int) (pagination.PaginatedResponse[Announcement], error)
	Create(announcement *Announcement) error
	Update(id string, announcement *Announcement) error
}

type announcementRepository struct {
	db *gorm.DB
}

func NewAnnouncementRepository() AnnouncementRepository {
	return announcementRepository{db: persistence.DB}
}

func (r announcementRepository) FindByID(id string, announcement *Announcement) error {
	return r.db.First(&announcement, "id = ?", id).Error
}

func (r announcementRepository) FindAll(itemsPerPage, pageNumber int) (pagination.PaginatedResponse[Announcement], error) {
	query := r.db.Select("ID", "Title", "Description", "CreatedAt")
	query = query.Find(&Announcement{}).Order("created_at DESC")
	return pagination.Find[Announcement](query, pageNumber, itemsPerPage)
}

func (r announcementRepository) Create(annoncement *Announcement) error {
	return r.db.Create(&annoncement).Error
}

func (r announcementRepository) Update(id string, announcement *Announcement) error {
	var foundAnnouncement Announcement

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
