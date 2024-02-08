package handlers

import (
	"funbanking/internal/domain/model"
	"funbanking/internal/domain/repository"
	"funbanking/internal/domain/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

type TransactionHandler struct {
	transactionService service.TransactionService
}

func NewTransactionHandler() TransactionHandler {
	return TransactionHandler{
		transactionService: service.NewTransactionService(
			repository.NewTransactionRepository(),
		),
	}
}

func (h TransactionHandler) FindByID(c *gin.Context) {
	id := c.Param("id")

	transaction, err := h.transactionService.FindByID(id)

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "Unable to find transaction"})
		return
	}

	c.JSON(http.StatusOK, transaction)
}

func (h TransactionHandler) Approve(c *gin.Context) {
	var transaction model.Transaction
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

func (h TransactionHandler) Decline(c *gin.Context) {
	var transaction model.Transaction
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

func (h TransactionHandler) Create(c *gin.Context) {
	var transaction model.Transaction

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
