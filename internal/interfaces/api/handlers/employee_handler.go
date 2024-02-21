package handlers

import (
	"funbanking/internal/domain/banking"
	"funbanking/internal/domain/users"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

type EmployeeHandler struct {
	employeeService banking.EmployeeService
	userService     users.UserService
}

func NewEmployeeHandler() EmployeeHandler {
	return EmployeeHandler{
		employeeService: banking.NewEmployeeService(
			banking.NewEmployeeRepository(),
		),
		userService: users.NewUserService(
			users.NewUserRepository(),
			nil,
			nil,
		),
	}
}

func (h EmployeeHandler) FindAllByUserID(c *gin.Context) {
	userID := c.Param("id")

	employees, err := h.employeeService.FindAllByUserID(userID)

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "Unable to find user"})
		return
	}

	c.JSON(http.StatusOK, employees)
}

func (h EmployeeHandler) FindAllByBankID(c *gin.Context) {
	employeeID := c.Param("id")

	employee, err := h.employeeService.FindAllByBankID(employeeID)

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "Unable to find bank"})
		return
	}

	c.JSON(http.StatusOK, employee)
}

func (h EmployeeHandler) Create(c *gin.Context) {
	var employee banking.Employee
	var request banking.NewEmployeeRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Malformed request"})
		return
	}

	user, err := h.userService.FindByUsernameOrEmail(request.Email)

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "User does not exist"})
		return
	}

	employee = banking.Employee{
		UserID: user.ID,
		BankID: request.BankID,
	}

	if err := h.employeeService.Create(&employee); err != nil {
		if strings.Contains(err.Error(), "duplicate") {
			c.JSON(http.StatusBadRequest, gin.H{"message": cases.Title(language.AmericanEnglish).String(user.FirstName) + " is already an employee at this bank"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Something went wrong"})
		return
	}

	c.JSON(http.StatusCreated, employee)
}

func (h EmployeeHandler) Delete(c *gin.Context) {
	id := c.Param("id")

	if err := h.employeeService.Delete(id); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "Unable to find customer"})
		return
	}

	c.JSON(http.StatusNoContent, nil)
}
