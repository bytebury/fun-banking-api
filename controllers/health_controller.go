package controllers

import (
	"golfer/models"
	"golfer/repositories"
	"golfer/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

type HealthController struct {
	healthService services.HealthService
	visitor       repositories.VisitorRepository
}

func NewHealthController(
	healthService services.HealthService,
	visitor repositories.VisitorRepository,
) *HealthController {
	return &HealthController{
		healthService,
		visitor,
	}
}

func (controller HealthController) GetHealthCheck(c *gin.Context) {
	var health models.Health

	if err := controller.healthService.GetHealthCheck(&health); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Unable to retrieve health information"})
		return
	}

	controller.visitor.AddVisitor(&models.Visitor{IPAddress: c.ClientIP()})

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

func (controller HealthController) GetVisitorInsights(c *gin.Context) {
	var result repositories.VisitorByDay

	if err := controller.visitor.GetVisitorsByDay(&result); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Something unexpected happened on the server"})
		return
	}

	c.JSON(http.StatusOK, result)
}
