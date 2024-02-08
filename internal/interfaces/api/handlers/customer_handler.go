package handlers

import (
	"funbanking/internal/domain/model"
	"funbanking/internal/domain/repository"
	"funbanking/internal/domain/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

type CustomerHandler struct {
	customerService service.CustomerService
}

func NewCustomerHandler() CustomerHandler {
	return CustomerHandler{
		customerService: service.NewCustomerService(
			repository.NewCustomerRepository(),
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
	var customer model.Customer

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
	var customer model.Customer
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
