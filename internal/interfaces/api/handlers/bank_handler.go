package handlers

import (
	"funbanking/internal/domain/repository"
	"funbanking/internal/domain/service"

	"github.com/gin-gonic/gin"
)

type BankHandler struct {
	bank service.BankService
}

func NewBankHandler() BankHandler {
	return BankHandler{
		bank: service.NewBankService(repository.NewBankRepository()),
	}
}

func (h BankHandler) FindByID(c *gin.Context) {
	h.bank.FindByID(c)
}

func (h BankHandler) FindAllByUserID(c *gin.Context) {
	h.bank.FindAllByUserID(c)
}

func (h BankHandler) FindByUsernameAndSlug(c *gin.Context) {
	h.bank.FindByUsernameAndSlug(c)
}

func (h BankHandler) Create(c *gin.Context) {
	h.bank.Create(c)
}

func (h BankHandler) Update(c *gin.Context) {
	h.bank.Update(c)
}

func (h BankHandler) Delete(c *gin.Context) {
	h.bank.Delete(c)
}
