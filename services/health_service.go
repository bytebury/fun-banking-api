package services

import (
	"golfer/models"
	"golfer/repositories"
)

func NewHealthService(repository repositories.HealthRepository) *HealthService {
	return &HealthService{repository}
}

type HealthService struct {
	repository repositories.HealthRepository
}

func (service HealthService) GetHealthCheck(health *models.Health) error {
	return service.repository.GetHealthCheck(health)
}

func (service HealthService) GetUserWeeklyInsights(insights *[]models.WeeklyInsights) error {
	return service.repository.GetUserWeeklyInsights(insights)
}
