package service

import (
	"funbanking/internal/domain/model"
	"funbanking/internal/domain/repository"
	"net/http"

	"github.com/gin-gonic/gin"
)

type TransactionService interface {
	FindByID(c *gin.Context)
	Approve(c *gin.Context)
	Decline(c *gin.Context)
	Create(c *gin.Context)
}

type transactionService struct {
	transactionRepository repository.TransactionRepository
}

func NewTransactionService(transactionRepository repository.TransactionRepository) TransactionService {
	return transactionService{transactionRepository}
}

func (s transactionService) FindByID(c *gin.Context) {
	var transaction model.Transaction
	transactionID := c.Param("id")

	if err := s.transactionRepository.FindByID(transactionID, &transaction); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Something went wrong"})
		return
	}

	c.JSON(http.StatusOK, transaction)
}

func (s transactionService) Approve(c *gin.Context) {
	var transaction model.Transaction
	transactionID := c.Param("id")

	if err := s.transactionRepository.FindByID(transactionID, &transaction); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Something went wrong"})
		return
	}

	if transaction.Status != model.TransactionPending {
		c.JSON(http.StatusBadRequest, gin.H{"message": "That transaction has already been processed"})
		return
	}

	transaction.Status = model.TransactionApproved

	if err := s.transactionRepository.Update(&transaction); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Something went wrong"})
		return
	}

	c.JSON(http.StatusAccepted, transaction)
}

func (s transactionService) Decline(c *gin.Context) {
	var transaction model.Transaction
	transactionID := c.Param("id")

	if err := s.transactionRepository.FindByID(transactionID, &transaction); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Something went wrong"})
		return
	}

	if transaction.Status != model.TransactionPending {
		c.JSON(http.StatusBadRequest, gin.H{"message": "That transaction has already been processed"})
		return
	}

	transaction.Status = model.TransactionDeclined

	if err := s.transactionRepository.Update(&transaction); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Something went wrong"})
		return
	}

	c.JSON(http.StatusAccepted, transaction)
}

func (s transactionService) Create(c *gin.Context) {
	var transaction model.Transaction

	if err := c.ShouldBindJSON(&transaction); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Malformed request"})
		return
	}

	if err := s.transactionRepository.Create(&transaction); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Something went wrong"})
		return
	}

	c.JSON(http.StatusCreated, transaction)
}
