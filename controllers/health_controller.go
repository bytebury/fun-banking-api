package controllers

import (
	"golfer/config"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetHealthCheckController(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"name":    config.AppName,
		"version": config.AppVersion,
		"message": "Everything is up and running!",
	})
}
