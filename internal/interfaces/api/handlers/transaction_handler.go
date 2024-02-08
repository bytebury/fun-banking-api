package handlers

import (
	"funbanking/internal/domain/repository"
	"funbanking/internal/domain/service"

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
	h.transactionService.FindByID(c)
}

func (h TransactionHandler) Approve(c *gin.Context) {
	h.transactionService.Approve(c)
}

func (h TransactionHandler) Decline(c *gin.Context) {
	h.transactionService.Decline(c)
}

func (h TransactionHandler) Create(c *gin.Context) {
	h.transactionService.Create(c)
}
