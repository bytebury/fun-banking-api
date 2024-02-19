package handlers

import (
	"funbanking/internal/domain/banking"
	"net/http"

	"github.com/gin-gonic/gin"
)

type TransactionHandler struct {
	transactionService banking.TransactionService
}

func NewTransactionHandler() TransactionHandler {
	return TransactionHandler{
		transactionService: banking.NewTransactionService(
			banking.NewTransactionRepository(),
		),
	}
}

// TODO: Similar to how notifications work today
func (h TransactionHandler) FindAllPendingTransactions(c *gin.Context) {
	c.JSON(http.StatusOK, nil)
}

// only employees can approve
func (h TransactionHandler) Approve(c *gin.Context) {
	var transaction banking.Transaction
	id := c.Param("id")

	if err := c.ShouldBindJSON(&transaction); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Malformed request"})
		return
	}

	if err := h.transactionService.Approve(id, &transaction); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Something went wrong"})
		return
	}

	c.JSON(http.StatusCreated, transaction)
}

// only employees can decline
func (h TransactionHandler) Decline(c *gin.Context) {
	var transaction banking.Transaction
	id := c.Param("id")

	if err := c.ShouldBindJSON(&transaction); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Malformed request"})
		return
	}

	if err := h.transactionService.Decline(id, &transaction); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Something went wrong"})
		return
	}

	c.JSON(http.StatusCreated, transaction)
}

// only customers part of the bank or employees can create
func (h TransactionHandler) Create(c *gin.Context) {
	var transaction banking.Transaction

	if err := c.ShouldBindJSON(&transaction); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Malformed request"})
		return
	}

	if err := h.transactionService.Create(&transaction); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Something went wrong"})
		return
	}

	c.JSON(http.StatusCreated, transaction)
}
