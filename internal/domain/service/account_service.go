package service

import (
	"funbanking/internal/domain/model"
	"funbanking/internal/domain/repository"
	"net/http"

	"github.com/gin-gonic/gin"
)

type AccountService interface {
	FindByID(ctx *gin.Context)
	Update(ctx *gin.Context)
}

type accountService struct {
	accountRepository repository.AccountRepository
}

func NewAccountService(accountRepository repository.AccountRepository) AccountService {
	return accountService{accountRepository}
}

func (s accountService) FindByID(c *gin.Context) {
	var account model.Account
	accountID := c.Param("id")

	if err := s.accountRepository.FindByID(accountID, &account); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Something went wrong"})
		return
	}

	c.JSON(http.StatusOK, account)
}

func (s accountService) Update(c *gin.Context) {
	var account model.Account

	if err := c.ShouldBindJSON(&account); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Malformed request"})
		return
	}

	if err := s.accountRepository.Update(&account); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Something went wrong"})
		return
	}

	c.JSON(http.StatusAccepted, account)
}
