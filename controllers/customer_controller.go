package controllers

import (
	"fmt"
	"golfer/config"
	"golfer/models"
	"golfer/services"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

type CustomerController struct {
	service         services.CustomerService
	bankService     services.BankService
	accountService  services.AccountService
	employeeService services.EmployeeService
	userService     services.UserService
}

func NewCustomerController(
	customer services.CustomerService,
	bankService services.BankService,
	accountService services.AccountService,
	employeeService services.EmployeeService,
	userService services.UserService,
) *CustomerController {
	return &CustomerController{customer, bankService, accountService, employeeService, userService}
}

// TODO: YOU SHOULDN'T BE ABLE TO LOOK UP CUSTOMERS IF THEY AREN'T IN A BANK YOU OWN!
func (controller CustomerController) FindByID(c *gin.Context) {
	customerID := c.Param("id")
	var customer models.Customer
	err := controller.service.FindByID(customerID, &customer)

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "customer not found"})
		return
	}

	c.JSON(http.StatusOK, customer)
}

func (controller CustomerController) FindAllAccounts(c *gin.Context) {
	customerID := c.Param("id")
	var accounts []models.Account

	err := controller.accountService.FindByCustomer(customerID, &accounts)

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "no accounts found"})
		return
	}

	if accounts == nil {
		accounts = make([]models.Account, 0)
	}

	c.JSON(http.StatusOK, accounts)
}

func (controller CustomerController) Create(c *gin.Context) {
	var customer models.Customer

	if err := c.ShouldBindJSON(&customer); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Bad request"})
		return
	}

	if !controller.isBankStaff(customer, c) {
		c.JSON(http.StatusForbidden, gin.H{"message": "You don't have access to that resource"})
		return
	}

	err := controller.service.Create(&customer)

	if isDuplicateError(err) {
		if strings.Contains(err.Error(), "pin") {
			c.JSON(http.StatusBadRequest, gin.H{"message": "A customer with that pin already exists"})
			return
		}
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Something went wrong creating your customer"})
		return
	}

	err = controller.accountService.Create(&models.Account{
		Name:       "Checkings",
		Balance:    float64(0),
		CustomerID: customer.ID,
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Something went wrong creating your account"})
		return
	}

	c.JSON(http.StatusCreated, customer)
}

func (controller CustomerController) Update(c *gin.Context) {
	var request models.Customer
	customerID := c.Param("id")

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Bad request"})
		return
	}

	var customer models.Customer
	if err := controller.service.FindByID(customerID, &customer); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "Customer does not exist"})
		return
	}

	if !controller.isBankStaff(customer, c) {
		c.JSON(http.StatusForbidden, gin.H{"message": "You don't have access to that resource"})
		return
	}

	customer, err := controller.service.Update(customerID, &request)

	if isDuplicateError(err) {
		if strings.Contains(err.Error(), "name") {
			c.JSON(http.StatusBadRequest, gin.H{"message": "That customer already exists"})
			return
		}

		if strings.Contains(err.Error(), "slug") {
			c.JSON(http.StatusBadRequest, gin.H{"message": "That customer already exists"})
			return
		}

		if strings.Contains(err.Error(), "pin") {
			c.JSON(http.StatusBadRequest, gin.H{"message": "A customer with that PIN already exists"})
			return
		}
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Something went wrong updating the customer"})
		return
	}

	c.JSON(http.StatusOK, customer)
}

func (controller CustomerController) Delete(c *gin.Context) {
	customerID := c.Param("id")

	var customer models.Customer
	if err := controller.service.FindByID(customerID, &customer); err != nil {
		c.JSON(http.StatusForbidden, gin.H{"message": "That customer doesn't exist"})
		return
	}

	if !controller.isBankStaff(customer, c) {
		c.JSON(http.StatusForbidden, gin.H{"message": "You don't have access to that resource"})
		return
	}

	if err := controller.service.Delete(customerID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Something went wrong deleting the customer"})
		return
	}

	c.JSON(http.StatusNoContent, nil)
}

func (controller CustomerController) Login(c *gin.Context) {
	var signInRequest models.CustomerSignInRequest

	if err := c.ShouldBindJSON(&signInRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Bad request"})
		return
	}

	var customer models.Customer
	if err := controller.service.Login(signInRequest.BankID, signInRequest.PIN, &customer); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Invalid login"})
		return
	}

	c.JSON(http.StatusOK, customer)
}

func (controller CustomerController) isBankStaff(customer models.Customer, c *gin.Context) bool {
	userID, err := strconv.Atoi(c.MustGet("user_id").(string))

	if err != nil {
		fmt.Println("⚠️ Unable to parse a userID when we should have been able to:", err)
		return false
	}

	var bank models.Bank
	if err := controller.bankService.FindByID(strconv.Itoa(int(customer.BankID)), &bank); err != nil {
		return false
	}

	if bank.UserID == uint(userID) {
		return true
	}

	var user models.User
	if err := controller.userService.FindByID(strconv.Itoa(userID), &user); err != nil {
		return false
	}

	if user.Role >= config.AdminRole {
		return true
	}

	var employees []models.Employee
	if err := controller.employeeService.FindByBank(strconv.Itoa(int(bank.ID)), &employees); err != nil {
		return false
	}

	for _, employee := range employees {
		if int(employee.UserID) == userID {
			return true
		}
	}

	return false
}
