package handlers

import (
	"funbanking/internal/domain/banking"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

type BankHandler struct {
	bankService banking.BankService
}

func NewBankHandler() BankHandler {
	return BankHandler{
		bankService: banking.NewBankService(
			banking.NewBankRepository(),
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

	if !h.bankService.IsOwner(bankID, userID) {
		c.JSON(http.StatusForbidden, gin.H{"message": "You do not have access to this bank"})
		return
	}

	c.JSON(http.StatusOK, bank)
}

func (h BankHandler) FindByUsernameAndSlug(c *gin.Context) {
	var request struct {
		Username string `json:"username"`
		Slug     string `json:"slug"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Malformed request"})
		return
	}

	bank, err := h.bankService.FindByUsernameAndSlug(request.Username, request.Slug)

	bank.User.Email = ""

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "Unable to find bank"})
		return
	}

	c.JSON(http.StatusOK, bank)
}

func (h BankHandler) FindAllByUserID(c *gin.Context) {
	userID := c.MustGet("user_id").(string)

	banks, err := h.bankService.FindAllByUserID(userID)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Something went wrong"})
		return
	}

	c.JSON(http.StatusOK, banks)
}

func (h BankHandler) FindAllCustomers(c *gin.Context) {
	bankID := c.Param("id")
	userID := c.MustGet("user_id").(string)

	customers, err := h.bankService.FindAllCustomers(bankID)

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "Unable to find bank"})
		return
	}

	if !h.bankService.IsOwner(bankID, userID) {
		c.JSON(http.StatusForbidden, gin.H{"message": "You do not have access to this bank"})
		return
	}

	c.JSON(http.StatusOK, customers)
}

func (h BankHandler) Create(c *gin.Context) {
	var bank banking.Bank

	if err := c.ShouldBindJSON(&bank); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Malformed request"})
		return
	}

	userID := c.MustGet("user_id").(string)

	if err := h.bankService.Create(userID, &bank); err != nil {
		if strings.Contains(err.Error(), "idx_user_slug") {
			c.JSON(http.StatusBadRequest, gin.H{"message": "You already have a bank by that name"})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{"message": "Something went wrong"})
		return
	}

	c.JSON(http.StatusCreated, bank)
}

func (h BankHandler) Update(c *gin.Context) {
	var bank banking.Bank
	bankID := c.Param("id")
	userID := c.MustGet("user_id").(string)

	if !h.bankService.IsOwner(bankID, userID) {
		c.JSON(http.StatusForbidden, gin.H{"message": "You do not have access to update this bank"})
		return
	}

	if err := c.ShouldBindJSON(&bank); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Malformed request"})
		return
	}

	if err := h.bankService.Update(bankID, &bank); err != nil {
		if strings.Contains(err.Error(), "idx_user_slug") {
			c.JSON(http.StatusBadRequest, gin.H{"message": "You already have a bank by that name"})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{"message": "Something went wrong"})
		return
	}

	c.JSON(http.StatusAccepted, bank)
}

func (h BankHandler) Delete(c *gin.Context) {
	bankID := c.Param("id")
	userID := c.MustGet("user_id").(string)

	if !h.bankService.IsOwner(bankID, userID) {
		c.JSON(http.StatusForbidden, gin.H{"message": "You do not have access to delete this bank"})
		return
	}

	if err := h.bankService.Delete(bankID); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "Unable to find bank"})
		return
	}

	c.JSON(http.StatusNoContent, nil)
}
