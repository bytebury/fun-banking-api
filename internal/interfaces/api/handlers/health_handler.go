package handlers

import (
	"funbanking/internal/domain/repository"
	"funbanking/internal/domain/service"

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
	h.health.GetHealthCheck(c)
}
