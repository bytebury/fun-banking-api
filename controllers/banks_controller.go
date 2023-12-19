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
	if err := controller.bank.FindByID(bankID, &bank); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "Bank not found"})
		return
	}

	if !controller.isBankOwner(c) {
		c.JSON(http.StatusForbidden, gin.H{"message": "You don't have access to that"})
		return
	}

	c.JSON(http.StatusOK, bank)
}

func (controller BankController) FindBanksByUserID(c *gin.Context) {
	var banks []models.Bank
	userID := c.MustGet("user_id").(string)
	err := controller.bank.FindBanksByUserID(userID, &banks)

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "No banks found"})
		return
	}

	if banks == nil {
		banks = make([]models.Bank, 0)
	}

	c.JSON(http.StatusOK, banks)
}

func (controller BankController) Create(c *gin.Context) {
	var bank models.Bank

	if err := c.ShouldBindJSON(&bank); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Bad request"})
		return
	}

	if bank.Name == "" || bank.Slug == "" || bank.Description == "" {
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

	if !controller.isBankOwner(c) {
		c.JSON(http.StatusForbidden, gin.H{"message": "You don't have access to that"})
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

	if !controller.isBankOwner(c) {
		c.JSON(http.StatusForbidden, gin.H{"message": "You don't have access to that"})
		return
	}

	if err := controller.bank.Delete(bankID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Something went wrong deleting the bank"})
		return
	}

	c.JSON(http.StatusNoContent, nil)
}

func (controller BankController) FindCustomers(c *gin.Context) {
	bankID := c.Param("id")
	var customers []models.Customer

	err := controller.bank.FindCustomers(bankID, &customers)

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "Error retrieving customers"})
		return
	}

	if customers == nil {
		customers = make([]models.Customer, 0)
	}

	c.JSON(http.StatusOK, customers)
}

/**
 * Determines if the current user can modify the given bank in context.
 */
func (controller BankController) isBankOwner(c *gin.Context) bool {
	bankID := c.Param("id")
	currentUserID := c.MustGet("user_id").(string)

	var bank models.Bank
	if err := controller.bank.FindByID(bankID, &bank); err != nil {
		return false
	}

	ownerID := strconv.Itoa(int(bank.User.ID))

	return ownerID == currentUserID
}
