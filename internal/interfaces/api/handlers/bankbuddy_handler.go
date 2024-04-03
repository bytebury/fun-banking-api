package handlers

import (
	"funbanking/internal/domain/banking"
	"net/http"
	"strings"

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
			banking.NewAccountService(
				banking.NewAccountRepository(),
				banking.NewCustomerRepository(),
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
		if strings.Contains(err.Error(), "insufficient funds") {
			c.JSON(http.StatusBadRequest, gin.H{"message": "You do not have sufficient funds to do that"})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{"message": "Something went wrong sending money to that customer"})
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
