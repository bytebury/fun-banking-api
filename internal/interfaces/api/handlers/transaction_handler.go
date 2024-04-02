package handlers

import (
	"funbanking/internal/domain/banking"
	"net/http"
	"strings"

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

func (h TransactionHandler) FindAllPendingTransactions(c *gin.Context) {
	userID := c.MustGet("user_id").(string)

	transactions, err := h.transactionService.FindAllPendingTransactions(userID)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Something went wrong"})
		return
	}

	c.JSON(http.StatusOK, transactions)
}

// only employees can approve
func (h TransactionHandler) Approve(c *gin.Context) {
	transactionID := c.Param("id")
	userID := c.MustGet("user_id").(string)

	transaction, err := h.transactionService.Approve(userID, transactionID)

	if err != nil {
		if strings.Contains(err.Error(), "processed") {
			c.JSON(http.StatusBadRequest, gin.H{"message": "This transaction has already been processed"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Something went wrong"})
		return
	}

	c.JSON(http.StatusAccepted, transaction)
}

// only employees can decline
func (h TransactionHandler) Decline(c *gin.Context) {
	transactionID := c.Param("id")
	userID := c.MustGet("user_id").(string)

	transaction, err := h.transactionService.Decline(userID, transactionID)

	if err != nil {
		if strings.Contains(err.Error(), "processed") {
			c.JSON(http.StatusBadRequest, gin.H{"message": "This transaction has already been processed"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Something went wrong"})
		return
	}

	c.JSON(http.StatusAccepted, transaction)
}

// only customers part of the bank or employees can create
func (h TransactionHandler) Create(c *gin.Context) {
	var transaction banking.Transaction

	userID := c.GetString("user_id")

	if err := c.ShouldBindJSON(&transaction); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Malformed request"})
		return
	}

	transaction.Status = "pending"

	if err := h.transactionService.Create(userID, &transaction); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Something went wrong"})
		return
	}

	c.JSON(http.StatusCreated, transaction)
}
