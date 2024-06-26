package handlers

import (
	"funbanking/internal/domain/announcements"
	"funbanking/package/utils"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type AnnouncementHandler struct {
	announcementService announcements.AnnouncementService
}

func NewAnnouncementHandler() AnnouncementHandler {
	return AnnouncementHandler{
		announcementService: announcements.NewAnnouncementService(
			announcements.NewAnnouncementRepository(),
		),
	}
}

func (h AnnouncementHandler) FindByID(c *gin.Context) {
	announcementID := c.Param("id")

	announcement, err := h.announcementService.FindByID(announcementID)

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "Unable to find announcement"})
		return
	}

	c.JSON(http.StatusOK, announcement)
}

func (h AnnouncementHandler) FindAll(c *gin.Context) {
	itemsPerPage, _ := strconv.Atoi(c.Query("itemsPerPage"))
	pageNumber, _ := strconv.Atoi(c.Query("pageNumber"))

	announcements, err := h.announcementService.FindAll(itemsPerPage, pageNumber)

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "Unable to find announcements"})
		return
	}

	c.JSON(http.StatusOK, announcements)
}

func (h AnnouncementHandler) Create(c *gin.Context) {
	userID := c.MustGet("user_id").(string)

	var announcement announcements.Announcement

	if err := c.ShouldBindJSON(&announcement); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Malformed request"})
		return
	}

	if userIDInt, err := utils.StringToUIntPointer(userID); err == nil {
		announcement.UserID = *userIDInt
	}

	if err := h.announcementService.Create(&announcement); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Something went wrong"})
		return
	}

	c.JSON(http.StatusCreated, announcement)
}

func (h AnnouncementHandler) Update(c *gin.Context) {
	userID := c.MustGet("user_id").(string)

	var announcement announcements.Announcement
	announcementID := c.Param("id")

	if err := c.ShouldBindJSON(&announcement); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Malformed request"})
		return
	}

	if userIDInt, err := utils.StringToUIntPointer(userID); err == nil {
		announcement.UserID = *userIDInt
	}

	if err := h.announcementService.Update(announcementID, &announcement); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Something went wrong"})
		return
	}

	c.JSON(http.StatusAccepted, announcement)
}
