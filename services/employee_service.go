package services

import (
	"errors"
	"golfer/models"
	"golfer/repositories"
	"strconv"
)

type EmployeeService struct {
	repository  repositories.EmployeeRepository
	userService UserService
	bankService BankService
}

func NewEmployeeService(
	repository repositories.EmployeeRepository,
	userService UserService,
	bankService BankService,
) *EmployeeService {
	return &EmployeeService{
		repository,
		userService,
		bankService,
	}
}

func (service EmployeeService) Create(request models.EmployeeRequest, employee *models.Employee, userID string) error {
	var user models.User
	if err := service.userService.FindByEmail(request.Email, &user); err != nil {
		return err
	}

	var bank models.Bank
	if err := service.bankService.FindByID(strconv.Itoa(int(request.BankID)), &bank); err != nil {
		return err
	}

	userIDNum, _ := strconv.Atoi(userID)
	if bank.UserID == uint(userIDNum) {
		return errors.New("you cannot add yourself as an employee")
	}

	employee.BankID = bank.ID
	employee.UserID = user.ID

	return service.repository.Create(employee)
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
