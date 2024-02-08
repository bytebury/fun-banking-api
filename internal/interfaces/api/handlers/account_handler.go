package handlers

import (
	"funbanking/internal/domain/model"
	"funbanking/internal/domain/repository"
	"funbanking/internal/domain/service"
	"funbanking/package/utils"
	"net/http"

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
	id := c.Param("id")

	account, err := h.accountService.FindByID(id)

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "Unable to find that account"})
		return
	}

	c.JSON(http.StatusOK, account)
}

func (h AccountHandler) FindTransactions(c *gin.Context) {
	id := c.Param("id")

	transactions, err := h.accountService.FindTransactions(id)

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "Unable to find that account"})
		return
	}

	c.JSON(http.StatusOK, utils.Listify(transactions))
}

func (h AccountHandler) Update(c *gin.Context) {
	var account model.Account
	id := c.Param("id")

	if err := c.ShouldBindJSON(&account); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Malformed request"})
		return
	}

	account, err := h.accountService.Update(id, &account)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Something went wrong"})
		return
	}

	c.JSON(http.StatusAccepted, account)
}
