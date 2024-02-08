package service

import (
	"funbanking/internal/domain/model"
	"funbanking/internal/domain/repository"
)

type HealthService interface {
	GetHealthCheck() (model.Health, error)
}

type healthService struct {
	healthRepository repository.HealthRepository
}

func NewHealthService(healthRepository repository.HealthRepository) HealthService {
	return healthService{healthRepository}
}

func (s healthService) GetHealthCheck() (model.Health, error) {
	var health model.Health
	err := s.healthRepository.GetHealthCheck(&health)
	return health, err
}
