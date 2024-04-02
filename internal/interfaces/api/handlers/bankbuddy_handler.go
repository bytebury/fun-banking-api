package handlers

import (
	"funbanking/internal/domain/banking"
	"net/http"

	"github.com/gin-gonic/gin"
)

type BankBuddyHandler struct {
	bankBuddyService banking.BankBuddyService
}

func NewBankBuddyHandler() BankBuddyHandler {
	return BankBuddyHandler{
		bankBuddyService: banking.NewBankBuddyService(
			banking.NewTransactionService(
				banking.NewTransactionRepository(),
			),
		),
	}
}

func (h BankBuddyHandler) Transfer(c *gin.Context) {
	var transferRequest banking.BankBuddyTransfer

	if err := c.ShouldBindJSON(&transferRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Malformed request"})
		return
	}

	if err := h.bankBuddyService.Transfer(&transferRequest); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Something went wrong"})
		return
	}

	c.JSON(http.StatusCreated, nil)
}

func (h BankBuddyHandler) FindRecipients(c *gin.Context) {
	customers, err := h.bankBuddyService.FindRecipients(c.Param("id"))

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Something went wrong"})
		return
	}

	c.JSON(http.StatusOK, customers)
}
