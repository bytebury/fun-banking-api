package handlers

import (
	"funbanking/internal/domain/banking"
	"net/http"

	"github.com/gin-gonic/gin"
)

type ConfigHandler struct{}

func NewConfigHandler() ConfigHandler {
	return ConfigHandler{}
}

func (h ConfigHandler) GetConfig(c *gin.Context) {
	c.JSON(http.StatusOK, banking.BankConfig)
}
