package controllers

import (
	"golfer/models"
	"golfer/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

type HealthController struct {
	healthService services.HealthService
}

func NewHealthController(
	healthService services.HealthService,
) *HealthController {
	return &HealthController{
		healthService,
	}
}

func (controller HealthController) GetHealthCheck(c *gin.Context) {
	var health models.Health

	if err := controller.healthService.GetHealthCheck(&health); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Unable to retrieve health information"})
		return
	}

	c.JSON(http.StatusOK, health)
}

func (controller HealthController) GetUserInsights(c *gin.Context) {
	var weeklyInsights []models.WeeklyInsights

	if err := controller.healthService.GetUserWeeklyInsights(&weeklyInsights); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Unable to retrieve user insights"})
		return
	}

	c.JSON(http.StatusOK, weeklyInsights)
}
