package service

import (
	"funbanking/internal/domain/model"
	"funbanking/internal/domain/repository"
)

type AnnouncementService interface {
	FindByID(id string) (model.Announcement, error)
	Create(announcement *model.Announcement) error
	Update(id string, announcement *model.Announcement) error
}

type announcementService struct {
	announcementRepository repository.AnnouncementRepository
}

func NewAnnouncementService(announcementRepository repository.AnnouncementRepository) AnnouncementService {
	return announcementService{announcementRepository}
}

func (s announcementService) FindByID(id string) (model.Announcement, error) {
	var announcement model.Announcement
	err := s.announcementRepository.FindByID(id, &announcement)
	return announcement, err
}

func (s announcementService) Create(announcement *model.Announcement) error {
	return s.announcementRepository.Create(announcement)
}

func (s announcementService) Update(id string, announcement *model.Announcement) error {
	return s.announcementRepository.Update(id, announcement)
}
