package controllers

import (
	"golfer/models"
	"golfer/services"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

type AnnouncementController struct {
	announcementService services.AnnouncementService
}

func NewAnnouncementController(
	announcementService services.AnnouncementService,
) *AnnouncementController {
	return &AnnouncementController{
		announcementService,
	}
}

func (controller AnnouncementController) Create(c *gin.Context) {
	var announcement models.Announcement
	if err := c.ShouldBindJSON(&announcement); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Bad request"})
		return
	}

	if err := controller.announcementService.Create(&announcement, c); err != nil {
		if strings.Contains(err.Error(), "missing required fields") {
			c.JSON(http.StatusBadRequest, gin.H{"message": "You are missing required fields"})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{"message": "Unable to create announcement"})
		return
	}

	c.JSON(http.StatusCreated, announcement)
}

func (controller AnnouncementController) Find(c *gin.Context) {
	var announcements []models.Announcement

	if err := controller.announcementService.Find(&announcements); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Something went wrong"})
		return
	}

	if announcements == nil {
		announcements = make([]models.Announcement, 0)
	}

	c.JSON(http.StatusOK, announcements)
}

func (controller AnnouncementController) FindByID(c *gin.Context) {
	announcementID := c.Param("id")

	var announcement models.Announcement
	if err := controller.announcementService.FindByID(announcementID, &announcement); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Unable to retrieve announcement"})
		return
	}

	c.JSON(http.StatusOK, announcement)
}

func (controller AnnouncementController) Update(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "This is coming soon."})
}

func (controller AnnouncementController) Delete(c *gin.Context) {
	announcementID := c.Param("id")

	if err := controller.announcementService.Delete(announcementID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Unable to delete announcement"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Deleted announcement"})
}
