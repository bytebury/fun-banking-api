package service

import (
	"funbanking/internal/domain/model"
	"funbanking/internal/domain/repository"
	"net/http"

	"github.com/gin-gonic/gin"
)

type CustomerService interface {
	FindByID(ctx *gin.Context)
	FindAccounts(ctx *gin.Context)
	Create(ctx *gin.Context)
	Update(ctx *gin.Context)
	Delete(ctx *gin.Context)
}

type customerService struct {
	customerRepository repository.CustomerRepository
}

func NewCustomerService(customerRepository repository.CustomerRepository) CustomerService {
	return customerService{customerRepository}
}

func (s customerService) FindByID(c *gin.Context) {
	var customer model.Customer
	customerID := c.Param("id")

	if err := s.customerRepository.FindByID(customerID, &customer); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Something went wrong"})
		return
	}

	c.JSON(http.StatusOK, customer)
}

func (s customerService) FindAccounts(c *gin.Context) {
	var accounts []model.Account
	customerID := c.Param("id")

	if err := s.customerRepository.FindAccounts(customerID, &accounts); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Something went wrong"})
		return
	}

	c.JSON(http.StatusOK, accounts)
}

func (s customerService) Create(c *gin.Context) {
	var customer model.Customer

	if err := c.ShouldBindJSON(&customer); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Malformed request"})
		return
	}

	if err := s.customerRepository.Create(&customer); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Something went wrong"})
		return
	}

	c.JSON(http.StatusCreated, customer)
}

func (s customerService) Update(c *gin.Context) {
	var customer model.Customer

	if err := c.ShouldBindJSON(&customer); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Malformed request"})
		return
	}

	if err := s.customerRepository.Update(&customer); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Something went wrong"})
		return
	}

	c.JSON(http.StatusAccepted, customer)
}

func (s customerService) Delete(c *gin.Context) {
	customerID := c.Param("id")

	if err := s.customerRepository.Delete(customerID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Something went wrong"})
		return
	}

	c.JSON(http.StatusNoContent, nil)
}
