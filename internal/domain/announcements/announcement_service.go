package announcements

import "funbanking/internal/infrastructure/pagination"

type AnnouncementService interface {
	FindByID(id string) (Announcement, error)
	FindAll(itemsPerPage, pageNumber int) (pagination.PaginatedResponse[Announcement], error)
	Create(announcement *Announcement) error
	Update(id string, announcement *Announcement) error
}

type announcementService struct {
	announcementRepository AnnouncementRepository
}

func NewAnnouncementService(announcementRepository AnnouncementRepository) AnnouncementService {
	return announcementService{announcementRepository}
}

func (s announcementService) FindByID(id string) (Announcement, error) {
	var announcement Announcement
	err := s.announcementRepository.FindByID(id, &announcement)
	return announcement, err
}

func (s announcementService) FindAll(itemsPerPage, pageNumber int) (pagination.PaginatedResponse[Announcement], error) {
	return s.announcementRepository.FindAll(itemsPerPage, pageNumber)
}

func (s announcementService) Create(announcement *Announcement) error {
	return s.announcementRepository.Create(announcement)
}

func (s announcementService) Update(id string, announcement *Announcement) error {
	return s.announcementRepository.Update(id, announcement)
}
