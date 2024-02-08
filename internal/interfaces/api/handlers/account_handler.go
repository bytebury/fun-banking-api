package handlers

import (
	"funbanking/internal/domain/repository"
	"funbanking/internal/domain/service"

	"github.com/gin-gonic/gin"
)

type AccountHandler struct {
	accountService service.AccountService
}

func NewAccountHandler() AccountHandler {
	return AccountHandler{
		accountService: service.NewAccountService(
			repository.NewAccountRepository(),
		),
	}
}

func (h AccountHandler) FindByID(c *gin.Context) {
	h.accountService.FindByID(c)
}

func (h AccountHandler) FindTransactions(c *gin.Context) {
	h.accountService.FindTransactions(c)
}

func (h AccountHandler) Update(c *gin.Context) {
	h.accountService.Update(c)
}
