package controllers

import (
	"golfer/models"
	"golfer/services"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

type BankController struct {
	bank services.BankService
}

func NewBankController(bank services.BankService) *BankController {
	return &BankController{bank}
}

func (controller BankController) FindByID(c *gin.Context) {
	bankID := c.Param("id")
	var bank models.Bank
	err := controller.bank.FindByID(bankID, &bank)

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "Bank not found"})
		return
	}

	c.JSON(http.StatusOK, bank)
}

func (controller BankController) Search(c *gin.Context) {
	var banks []models.Bank
	err := controller.bank.Search(c, &banks)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Something went wrong"})
		return
	}

	c.JSON(http.StatusOK, banks)
}

func (controller BankController) Create(c *gin.Context) {
	var bank models.Bank

	if err := c.ShouldBindJSON(&bank); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Bad request"})
		return
	}

	currentUserID, userErr := strconv.Atoi(c.MustGet("user_id").(string))

	if userErr != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid user"})
		return
	}

	bank.UserID = uint(currentUserID)

	err := controller.bank.Create(&bank)

	if isDuplicateError(err) {
		if strings.Contains(err.Error(), "name") {
			c.JSON(http.StatusBadRequest, gin.H{"message": "That bank already exists"})
			return
		}

		if strings.Contains(err.Error(), "slug") {
			c.JSON(http.StatusBadRequest, gin.H{"message": "That bank already exists"})
			return
		}
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Something went wrong creating your bank"})
		return
	}

	c.JSON(http.StatusCreated, bank)
}

func (controller BankController) Update(c *gin.Context) {
	var request models.Bank
	bankID := c.Param("id")

	if !controller.canModify(c) {
		c.JSON(http.StatusForbidden, gin.H{"message": "You don't have access to that resource"})
		return
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Bad request"})
		return
	}

	bank, err := controller.bank.Update(bankID, &request)

	if isDuplicateError(err) {
		if strings.Contains(err.Error(), "name") {
			c.JSON(http.StatusBadRequest, gin.H{"message": "That bank already exists"})
			return
		}

		if strings.Contains(err.Error(), "slug") {
			c.JSON(http.StatusBadRequest, gin.H{"message": "That bank already exists"})
			return
		}
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Something went wrong updating the bank"})
		return
	}

	c.JSON(http.StatusOK, bank)
}

func (controller BankController) Delete(c *gin.Context) {
	bankID := c.Param("id")

	if !controller.canModify(c) {
		c.JSON(http.StatusForbidden, gin.H{"message": "You don't have access to that resource"})
		return
	}

	if err := controller.bank.Delete(bankID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Something went wrong deleting the bank"})
		return
	}

	c.JSON(http.StatusNoContent, nil)
}

/**
 * Determines if the current user can modify the given bank in context.
 */
func (controller BankController) canModify(c *gin.Context) bool {
	bankID := c.Param("id")
	var bank models.Bank

	err := controller.bank.FindByID(bankID, &bank)

	if err != nil {
		return false
	}

	ownerId := strconv.Itoa(int(bank.User.ID))
	currentUserID := c.MustGet("user_id").(string)
	return ownerId == currentUserID
}
