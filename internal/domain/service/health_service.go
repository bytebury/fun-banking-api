package service

import (
	"funbanking/internal/domain/model"
	"funbanking/internal/domain/repository"
	"net/http"

	"github.com/gin-gonic/gin"
)

type HealthService struct {
	repo repository.Health
}

func NewHealthService(healthRepository repository.Health) HealthService {
	return HealthService{repo: healthRepository}
}

func (s HealthService) GetHealthCheck(c *gin.Context) {
	var health model.Health
	err := s.repo.GetHealthCheck(&health)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Something went wrong"})
		return
	}

	c.JSON(http.StatusOK, health)
}
