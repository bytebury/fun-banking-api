package controllers

import (
	"golfer/config"
	"golfer/models"
	"golfer/services"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

type BankController struct {
	bank            services.BankService
	userService     services.UserService
	employeeService services.EmployeeService
}

func NewBankController(bank services.BankService, userService services.UserService, employeeService services.EmployeeService) *BankController {
	return &BankController{bank, userService, employeeService}
}

func (controller BankController) FindByID(c *gin.Context) {
	bankID := c.Param("id")

	var bank models.Bank
	if err := controller.bank.FindByID(bankID, &bank); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "Bank not found"})
		return
	}

	if !controller.isBankOwner(c) && !controller.isBankEmployee(c) {
		c.JSON(http.StatusForbidden, gin.H{"message": "You don't have access to that"})
		return
	}

	c.JSON(http.StatusOK, bank)
}

func (controller BankController) FindBanksByUserID(c *gin.Context) {
	var banks []models.Bank
	var employeeOf []models.Employee
	var employeeBanks []models.Bank = make([]models.Bank, 0)

	userID := c.MustGet("user_id").(string)

	if err := controller.bank.FindBanksByUserID(userID, &banks); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "Unable to retrieve bank information"})
		return
	}

	if err := controller.employeeService.FindByUser(userID, &employeeOf); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "Unable to retrieve bank information"})
		return
	}

	for _, employee := range employeeOf {
		employeeBanks = append(employeeBanks, employee.Bank)
	}

	if banks == nil {
		banks = make([]models.Bank, 0)
	}

	banks = append(banks, employeeBanks...)

	c.JSON(http.StatusOK, banks)
}

func (controller BankController) FindByUsernameAndSlug(c *gin.Context) {
	var bank models.Bank

	username := c.Param("username")
	slug := c.Param("slug")

	if err := controller.bank.FindByUsernameAndSlug(username, slug, &bank); err != nil {
		if strings.Contains(err.Error(), "record not found") {
			c.JSON(http.StatusNotFound, gin.H{"message": "Unable to retrieve that bank"})
			return
		}

		c.JSON(http.StatusBadRequest, gin.H{"message": "Unable to retrieve that bank"})
		return
	}

	c.JSON(http.StatusOK, bank)
}

func (controller BankController) Create(c *gin.Context) {
	var bank models.Bank

	if err := c.ShouldBindJSON(&bank); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Bad request"})
		return
	}

	if bank.Name == "" || bank.Slug == "" || bank.Description == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Bad request"})
		return
	}

	currentUserID, userErr := strconv.Atoi(c.MustGet("user_id").(string))

	if userErr != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid user"})
		return
	}

	bank.UserID = uint(currentUserID)

	err := controller.bank.Create(&bank)

	if isDuplicateError(err) {
		if strings.Contains(err.Error(), "name") {
			c.JSON(http.StatusBadRequest, gin.H{"message": "That bank already exists"})
			return
		}

		if strings.Contains(err.Error(), "slug") {
			c.JSON(http.StatusBadRequest, gin.H{"message": "That bank already exists"})
			return
		}
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Something went wrong creating your bank"})
		return
	}

	c.JSON(http.StatusCreated, bank)
}

func (controller BankController) Update(c *gin.Context) {
	var request models.Bank
	bankID := c.Param("id")

	if !controller.isBankOwner(c) {
		c.JSON(http.StatusForbidden, gin.H{"message": "You don't have access to that"})
		return
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Bad request"})
		return
	}

	bank, err := controller.bank.Update(bankID, &request)

	if isDuplicateError(err) {
		if strings.Contains(err.Error(), "name") {
			c.JSON(http.StatusBadRequest, gin.H{"message": "That bank already exists"})
			return
		}

		if strings.Contains(err.Error(), "slug") {
			c.JSON(http.StatusBadRequest, gin.H{"message": "That bank already exists"})
			return
		}
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Something went wrong updating the bank"})
		return
	}

	c.JSON(http.StatusOK, bank)
}

func (controller BankController) Delete(c *gin.Context) {
	bankID := c.Param("id")

	if !controller.isBankOwner(c) {
		c.JSON(http.StatusForbidden, gin.H{"message": "You don't have access to that"})
		return
	}

	if err := controller.bank.Delete(bankID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Something went wrong deleting the bank"})
		return
	}

	c.JSON(http.StatusNoContent, nil)
}

func (controller BankController) FindCustomers(c *gin.Context) {
	bankID := c.Param("id")
	var customers []models.Customer

	err := controller.bank.FindCustomers(bankID, &customers)

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "Error retrieving customers"})
		return
	}

	if customers == nil {
		customers = make([]models.Customer, 0)
	}

	c.JSON(http.StatusOK, customers)
}

func (controller BankController) isBankEmployee(c *gin.Context) bool {
	bankID := c.Param("id")
	currentUserID := c.MustGet("user_id").(string)

	var employees []models.Employee
	if err := controller.employeeService.FindByBank(bankID, &employees); err != nil {
		return false
	}

	for _, employee := range employees {
		if strconv.Itoa(int(employee.UserID)) == currentUserID {
			return true
		}
	}

	return false
}

/**
 * Determines if the current user can modify the given bank in context.
 */
func (controller BankController) isBankOwner(c *gin.Context) bool {
	bankID := c.Param("id")
	currentUserID := c.MustGet("user_id").(string)

	var user models.User
	if err := controller.userService.FindByID(currentUserID, &user); err != nil {
		return false
	}

	// Admins can access everything
	if user.Role >= config.AdminRole {
		return true
	}

	var bank models.Bank
	if err := controller.bank.FindByID(bankID, &bank); err != nil {
		return false
	}

	ownerID := strconv.Itoa(int(bank.User.ID))

	return ownerID == currentUserID
}
