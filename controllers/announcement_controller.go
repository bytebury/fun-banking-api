package controllers

import (
	"golfer/models"
	"golfer/services"
	"net/http"
	"strconv"
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
	limit, limitErr := strconv.Atoi(c.Query("limit"))
	page, pageErr := strconv.Atoi(c.Query("page"))

	if limitErr != nil || limit <= 0 {
		limit = 5
	}

	if pageErr != nil || page <= 0 {
		page = 1
	}

	announcements, err := controller.announcementService.Find(limit, page)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Something went wrong"})
		return
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
	announcementID := c.Param("id")

	var updateAnnouncementRequest models.UpdateAnnouncementRequest
	if err := c.ShouldBindJSON(&updateAnnouncementRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Bad request"})
		return
	}

	if err := controller.announcementService.Update(announcementID, &updateAnnouncementRequest, c); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Unable to update the announcemnet"})
		return
	}

	c.JSON(http.StatusAccepted, gin.H{"message": "Announcement updated"})
}

func (controller AnnouncementController) Delete(c *gin.Context) {
	announcementID := c.Param("id")

	if err := controller.announcementService.Delete(announcementID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Unable to delete announcement"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Deleted announcement"})
}
