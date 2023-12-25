package controllers

import (
	"golfer/models"
	"golfer/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

type EmployeeController struct {
	employeeService services.EmployeeService
}

func NewEmployeeController(
	employeeService services.EmployeeService,
) *EmployeeController {
	return &EmployeeController{
		employeeService,
	}
}

func (controller EmployeeController) Create(c *gin.Context) {
	var employee models.Employee

	if err := c.ShouldBindJSON(&employee); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Bad request"})
		return
	}

	if err := controller.employeeService.Create(&employee); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Could not add user as an employee"})
		return
	}

	c.JSON(http.StatusOK, employee)
}

func (controller EmployeeController) FindByBank(c *gin.Context) {
	bankID := c.Param("id")

	var employees []models.Employee
	if err := controller.employeeService.FindByBank(bankID, &employees); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "Unable to find and employees for that bank"})
		return
	}

	c.JSON(http.StatusOK, employees)
}

func (controller EmployeeController) FindByUser(c *gin.Context) {
	userID := c.Param("id")

	var employees []models.Employee
	if err := controller.employeeService.FindByUser(userID, &employees); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "Unable to find employee data for that user"})
		return
	}

	c.JSON(http.StatusOK, employees)
}

func (controller EmployeeController) Delete(c *gin.Context) {
	employeeID := c.Param("id")

	if err := controller.employeeService.Delete(employeeID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Unable to delete that employee"})
		return
	}

	c.JSON(http.StatusNoContent, nil)
}
