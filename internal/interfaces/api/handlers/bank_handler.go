package handlers

import (
	"funbanking/internal/domain/model"
	"funbanking/internal/domain/repository"
	"funbanking/internal/domain/service"
	"net/http"

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
	bankID := c.Param("id")
	userID := c.MustGet("user_id").(string)

	bank, err := h.bankService.FindByID(bankID)

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "Unable to find bank"})
		return
	}

	if !h.bankService.IsBankOwner(bankID, userID) {
		c.JSON(http.StatusForbidden, gin.H{"message": "You do not have access to this bank"})
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

	bank, err := h.bankService.FindByUsernameAndSlug(request.username, request.slug)

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "Unable to find bank"})
		return
	}

	c.JSON(http.StatusOK, bank)
}

func (h BankHandler) FindAllCustomers(c *gin.Context) {
	bankID := c.Param("id")
	userID := c.MustGet("user_id").(string)

	customers, err := h.bankService.FindAllCustomers(bankID)

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "Unable to find bank"})
		return
	}

	if !h.bankService.IsBankOwner(bankID, userID) {
		c.JSON(http.StatusForbidden, gin.H{"message": "You do not have access to this bank"})
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
	bankID := c.Param("id")
	userID := c.MustGet("user_id").(string)

	if err := c.ShouldBindJSON(&bank); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Malformed request"})
		return
	}

	if err := h.bankService.Update(bankID, &bank); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Something went wrong"})
		return
	}

	if !h.bankService.IsBankOwner(bankID, userID) {
		c.JSON(http.StatusForbidden, gin.H{"message": "You do not have access to update this bank"})
		return
	}

	c.JSON(http.StatusAccepted, bank)
}

func (h BankHandler) Delete(c *gin.Context) {
	bankID := c.Param("id")
	userID := c.MustGet("user_id").(string)

	if err := h.bankService.Delete(bankID); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "Unable to find bank"})
		return
	}

	if !h.bankService.IsBankOwner(bankID, userID) {
		c.JSON(http.StatusForbidden, gin.H{"message": "You do not have access to delete this bank"})
		return
	}

	c.JSON(http.StatusNoContent, nil)
}
