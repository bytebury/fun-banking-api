package handlers

import (
	"funbanking/internal/domain/model"
	"funbanking/internal/domain/repository"
	"funbanking/internal/domain/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

type AnnouncementHandler struct {
	announcementService service.AnnouncementService
}

func NewAnnouncementHandler() AnnouncementHandler {
	return AnnouncementHandler{
		announcementService: service.NewAnnouncementService(
			repository.NewAnnouncementRepository(),
		),
	}
}

func (h AnnouncementHandler) FindByID(c *gin.Context) {
	id := c.Param("id")

	announcement, err := h.announcementService.FindByID(id)

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "Unable to find announcement"})
		return
	}

	c.JSON(http.StatusOK, announcement)
}

func (h AnnouncementHandler) Create(c *gin.Context) {
	var announcement model.Announcement

	if err := c.ShouldBindJSON(&announcement); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Malformed request"})
		return
	}

	if err := h.announcementService.Create(&announcement); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Something went wrong"})
		return
	}

	c.JSON(http.StatusCreated, announcement)
}

func (h AnnouncementHandler) Update(c *gin.Context) {
	var announcement model.Announcement
	id := c.Param("id")

	if err := c.ShouldBindJSON(&announcement); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Malformed request"})
		return
	}

	if err := h.announcementService.Update(id, &announcement); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Something went wrong"})
		return
	}

	c.JSON(http.StatusAccepted, announcement)
}
