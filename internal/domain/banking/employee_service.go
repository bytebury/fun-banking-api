package banking

import (
	"errors"
	"funbanking/package/utils"
	"strconv"
)

type EmployeeService interface {
	FindAllByBankID(bankID string) ([]Employee, error)
	FindAllByUserID(userID string) ([]Employee, error)
	Create(employee *Employee) error
	Delete(employeeID string) error
}

type employeeService struct {
	employeeRepository EmployeeRepository
}

func NewEmployeeService(employeeRepository EmployeeRepository) EmployeeService {
	return employeeService{employeeRepository}
}

func (s employeeService) FindAllByBankID(bankID string) ([]Employee, error) {
	var employees []Employee
	err := s.employeeRepository.FindAllByBankID(bankID, &employees)
	return utils.Listify(employees), err
}

func (s employeeService) FindAllByUserID(userID string) ([]Employee, error) {
	var employees []Employee
	err := s.employeeRepository.FindAllByUserID(userID, &employees)
	return utils.Listify(employees), err
}

func (s employeeService) Create(employee *Employee) error {
	if employees, err := s.FindAllByBankID(strconv.Itoa(int(employee.BankID))); err != nil {
		return err
	} else if EnablePremium && len(employees) >= BankConfig.Limits.Free.Employees {
		return errors.New("limit reached")
	}

	return s.employeeRepository.Create(employee)
}

func (s employeeService) Delete(employeeID string) error {
	return s.employeeRepository.Delete(employeeID)
}
