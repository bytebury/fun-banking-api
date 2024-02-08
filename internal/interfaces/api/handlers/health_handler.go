package handlers

import (
	"funbanking/internal/domain/repository"
	"funbanking/internal/domain/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

type HealthHandler struct {
	health service.HealthService
}

func NewHealthHandler() HealthHandler {
	return HealthHandler{
		health: service.NewHealthService(repository.NewHealthRepository()),
	}
}

func (h HealthHandler) GetHealthCheck(c *gin.Context) {
	healthMetrics, err := h.health.GetHealthCheck()

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Something went wrong"})
		return
	}

	c.JSON(http.StatusOK, healthMetrics)
}
