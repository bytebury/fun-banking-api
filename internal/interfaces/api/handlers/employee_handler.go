package handlers

import (
	"funbanking/internal/domain/banking"
	"net/http"

	"github.com/gin-gonic/gin"
)

type EmployeeHandler struct {
	employeeService banking.EmployeeService
}

func NewEmployeeHandler() EmployeeHandler {
	return EmployeeHandler{
		employeeService: banking.NewEmployeeService(
			banking.NewEmployeeRepository(),
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

	if err := c.ShouldBindJSON(&employee); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Malformed request"})
		return
	}

	if err := h.employeeService.Create(&employee); err != nil {
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
