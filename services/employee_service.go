package services

import (
	"golfer/models"
	"golfer/repositories"
)

type EmployeeService struct {
	repository repositories.EmployeeRepository
}

func NewEmployeeService(repository repositories.EmployeeRepository) *EmployeeService {
	return &EmployeeService{
		repository,
	}
}

func (service EmployeeService) Create(request *models.Employee) error {
	return service.repository.Create(request)
}

func (service EmployeeService) FindByBank(bankID string, employees *[]models.Employee) error {
	return service.repository.FindByBank(bankID, employees)
}

func (service EmployeeService) FindByUser(userID string, employees *[]models.Employee) error {
	return service.repository.FindByUser(userID, employees)
}

func (service EmployeeService) Delete(employeeID string) error {
	return service.repository.Delete(employeeID)
}
