package service

import (
	"funbanking/internal/domain/model"
	"funbanking/internal/domain/repository"
	"net/http"

	"github.com/gin-gonic/gin"
)

type HealthService struct {
	repo repository.HealthRepository
}

func NewHealthService(healthRepository repository.HealthRepository) HealthService {
	return HealthService{repo: healthRepository}
}

func (s HealthService) GetHealthCheck(c *gin.Context) {
	var health model.Health

	if err := s.repo.GetHealthCheck(&health); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Something went wrong"})
		return
	}

	c.JSON(http.StatusOK, health)
}
