package services

import (
	"errors"
	"golfer/models"
	"golfer/repositories"
	"strconv"

	"github.com/gin-gonic/gin"
)

func NewAnnouncementService(repository repositories.AnnouncementRepository) *AnnouncementService {
	return &AnnouncementService{repository}
}

type AnnouncementService struct {
	repository repositories.AnnouncementRepository
}

func (service AnnouncementService) Create(announcement *models.Announcement, c *gin.Context) error {
	announcement.UserID = service.getUserID(c)

	if announcement.Title == "" || announcement.Description == "" {
		return errors.New("you are missing required fields")
	}

	return service.repository.Create(announcement)
}

func (service AnnouncementService) Find(announcements *[]models.Announcement) error {
	return service.repository.Find(announcements)
}

func (service AnnouncementService) FindByID(announcementID string, announcement *models.Announcement) error {
	return service.repository.FindByID(announcementID, announcement)
}

func (service AnnouncementService) Update(announcementID string, request *models.UpdateAnnouncementRequest, c *gin.Context) error {
	var announcement models.Announcement
	if err := service.repository.FindByID(announcementID, &announcement); err != nil {
		return err
	}

	announcement.UserID = service.getUserID(c)

	if request.Title != "" {
		announcement.Title = request.Title
	}

	if request.Description != "" {
		announcement.Description = request.Description
	}

	return service.repository.Update(&announcement)
}

func (service AnnouncementService) Delete(announcementID string) error {
	return service.repository.Delete(announcementID)
}

func (service AnnouncementService) getUserID(c *gin.Context) uint {
	userIDStr := c.MustGet("user_id").(string)
	userID, _ := strconv.ParseUint(userIDStr, 10, 64)
	return uint(userID)
}
