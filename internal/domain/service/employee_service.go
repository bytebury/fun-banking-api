package service

import (
	"funbanking/internal/domain/model"
	"funbanking/internal/domain/repository"
	"funbanking/package/utils"
)

type EmployeeService interface {
	FindAllByBankID(bankID string) ([]model.Employee, error)
	FindAllByUserID(userID string) ([]model.Employee, error)
	Create(employee *model.Employee) error
	Delete(id string) error
}

type employeeService struct {
	employeeRepository repository.EmployeeRepository
}

func NewEmployeeService(employeeRepository repository.EmployeeRepository) EmployeeService {
	return employeeService{employeeRepository}
}

func (s employeeService) FindAllByBankID(bankID string) ([]model.Employee, error) {
	var employees []model.Employee
	err := s.employeeRepository.FindAllByBankID(bankID, &employees)
	return utils.Listify(employees), err
}

func (s employeeService) FindAllByUserID(userID string) ([]model.Employee, error) {
	var employees []model.Employee
	err := s.employeeRepository.FindAllByUserID(userID, &employees)
	return utils.Listify(employees), err
}

func (s employeeService) Create(employee *model.Employee) error {
	return s.employeeRepository.Create(employee)
}

func (s employeeService) Delete(id string) error {
	return s.employeeRepository.Delete(id)
}
