package handlers

import (
	"funbanking/internal/domain/model"
	"funbanking/internal/domain/repository"
	"funbanking/internal/domain/service"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

type BankHandler struct {
	bankService service.BankService
}

func NewBankHandler() BankHandler {
	return BankHandler{
		bankService: service.NewBankService(
			repository.NewBankRepository(),
		),
	}
}

func (h BankHandler) FindByID(c *gin.Context) {
	id := c.Param("id")

	bank, err := h.bankService.FindByID(id)

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "Unable to find bank"})
		return
	}

	c.JSON(http.StatusOK, bank)
}

func (h BankHandler) FindByUsernameAndSlug(c *gin.Context) {
	var request struct {
		username string
		slug     string
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Malformed request"})
		return
	}

	bank, err := h.bankService.FindByUsernameAndSlug(
		strings.TrimSpace(request.username),
		strings.TrimSpace(request.slug),
	)

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "Unable to find bank"})
		return
	}

	c.JSON(http.StatusOK, bank)
}

func (h BankHandler) FindAllCustomers(c *gin.Context) {
	id := c.Param("id")

	customers, err := h.bankService.FindAllCustomers(id)

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "Unable to find bank"})
		return
	}

	c.JSON(http.StatusOK, customers)
}

func (h BankHandler) Create(c *gin.Context) {
	var bank model.Bank

	if err := c.ShouldBindJSON(&bank); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Malformed request"})
		return
	}

	if err := h.bankService.Create(&bank); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Something went wrong"})
		return
	}

	c.JSON(http.StatusCreated, bank)
}

func (h BankHandler) Update(c *gin.Context) {
	var bank model.Bank
	id := c.Param("id")

	if err := c.ShouldBindJSON(&bank); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Malformed request"})
		return
	}

	if err := h.bankService.Update(id, &bank); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Something went wrong"})
		return
	}

	c.JSON(http.StatusAccepted, bank)
}

func (h BankHandler) Delete(c *gin.Context) {
	id := c.Param("id")

	if err := h.bankService.Delete(id); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "Unable to find bank"})
		return
	}

	c.JSON(http.StatusNoContent, nil)
}
