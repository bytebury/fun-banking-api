package handlers

import (
	"funbanking/internal/domain/banking"
	"funbanking/package/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

type AccountHandler struct {
	accountService banking.AccountService
}

func NewAccountHandler() AccountHandler {
	return AccountHandler{
		accountService: banking.NewAccountService(
			banking.NewAccountRepository(),
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
	var account banking.Account
	id := c.Param("id")

	if err := c.ShouldBindJSON(&account); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Malformed request"})
		return
	}

	if err := h.accountService.Update(id, &account); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Something went wrong"})
		return
	}

	c.JSON(http.StatusAccepted, account)
}
