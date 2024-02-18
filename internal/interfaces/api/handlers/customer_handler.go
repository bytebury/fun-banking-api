package handlers

import (
	"funbanking/internal/domain/banking"
	"funbanking/internal/infrastructure/auth"
	"net/http"

	"github.com/gin-gonic/gin"
)

type CustomerHandler struct {
	customerService banking.CustomerService
}

func NewCustomerHandler() CustomerHandler {
	customerRepository := banking.NewCustomerRepository()
	return CustomerHandler{
		customerService: banking.NewCustomerService(
			customerRepository,
			auth.NewCustomerAuth(customerRepository),
		),
	}
}

func (h CustomerHandler) FindByID(c *gin.Context) {
	id := c.Param("id")

	customer, err := h.customerService.FindByID(id)

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "Unable to find customer"})
		return
	}

	c.JSON(http.StatusOK, customer)
}

func (h CustomerHandler) FindAccounts(c *gin.Context) {
	id := c.Param("id")

	accounts, err := h.customerService.FindAccounts(id)

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "Unable to find customer"})
		return
	}

	c.JSON(http.StatusOK, accounts)
}

func (h CustomerHandler) Create(c *gin.Context) {
	var customer banking.Customer

	if err := c.ShouldBindJSON(&customer); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Malformed request"})
		return
	}

	if err := h.customerService.Create(&customer); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Something went wrong"})
		return
	}

	c.JSON(http.StatusCreated, customer)
}

func (h CustomerHandler) Update(c *gin.Context) {
	var customer banking.Customer
	id := c.Param("id")

	if err := c.ShouldBindJSON(&customer); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Malformed request"})
		return
	}

	if err := h.customerService.Update(id, &customer); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Something went wrong"})
		return
	}

	c.JSON(http.StatusCreated, customer)
}

func (h CustomerHandler) Delete(c *gin.Context) {
	id := c.Param("id")

	if err := h.customerService.Delete(id); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "Unable to find customer"})
		return
	}

	c.JSON(http.StatusNoContent, nil)
}

func (h CustomerHandler) Login(c *gin.Context) {
	var request banking.LoginRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Malformed request"})
		return
	}

	token, customer, err := h.customerService.Login(request.BankID, request.PIN)

	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Unable to log you in, invalid credentials"})
		return
	}

	response := struct {
		Token    string           `json:"token"`
		Customer banking.Customer `json:"customer"`
	}{Token: token, Customer: customer}

	c.JSON(http.StatusOK, response)
}
