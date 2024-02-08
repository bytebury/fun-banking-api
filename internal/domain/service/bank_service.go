package service

import (
	"funbanking/internal/domain/model"
	"funbanking/internal/domain/repository"
	"funbanking/package/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

type BankService interface {
	FindByID(ctx *gin.Context)
	FindAllMyBanks(ctx *gin.Context)
	FindByUsernameAndSlug(ctx *gin.Context)
	FindAllCustomers(ctx *gin.Context)
	Create(ctx *gin.Context)
	Update(ctx *gin.Context)
	Delete(ctx *gin.Context)
}

type bankService struct {
	bankRepository repository.BankRepository
}

func NewBankService(bankRepository repository.BankRepository) BankService {
	return bankService{bankRepository}
}

func (s bankService) FindByID(c *gin.Context) {
	var bank model.Bank
	bankID := c.Param("id")

	if err := s.bankRepository.FindByID(bankID, &bank); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Something went wrong"})
		return
	}

	c.JSON(http.StatusOK, bank)
}

func (s bankService) FindAllMyBanks(c *gin.Context) {
	var banks []model.Bank
	// TODO: This will come from context from the token!
	userID := c.Param("user_id")

	if err := s.bankRepository.FindAllByUserID(userID, &banks); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Something went wrong"})
		return
	}

	c.JSON(http.StatusOK, utils.Listify[model.Bank](banks))
}

func (s bankService) FindByUsernameAndSlug(c *gin.Context) {
	var bank model.Bank

	var request struct {
		username string
		slug     string
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Malformed request"})
		return
	}

	if err := s.bankRepository.FindByUsernameAndSlug(request.username, request.slug, &bank); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Something went wrong"})
		return
	}

	c.JSON(http.StatusOK, bank)
}

// TODO: THIS SHOULD BE PAGINATED
func (s bankService) FindAllCustomers(c *gin.Context) {
	var customers []model.Customer
	bankID := c.Param("id")

	if err := s.bankRepository.FindAllCustomers(bankID, &customers); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Something went wrong"})
		return
	}

	c.JSON(http.StatusOK, utils.Listify(customers))
}

func (s bankService) Create(c *gin.Context) {
	var bank model.Bank

	if err := c.ShouldBindJSON(&bank); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Malformed request"})
		return
	}

	// Do validations

	if err := s.bankRepository.Create(&bank); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Something went wrong"})
		return
	}

	c.JSON(http.StatusCreated, bank)
}

func (s bankService) Update(c *gin.Context) {
	var bank model.Bank
	bankID := c.Param("id")

	if err := c.ShouldBindJSON(&bank); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Malformed request"})
		return
	}

	// Do validations

	if err := s.bankRepository.Update(bankID, &bank); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Something went wrong"})
		return
	}

	c.JSON(http.StatusAccepted, bank)
}

func (s bankService) Delete(c *gin.Context) {
	bankID := c.Param("id")

	if err := s.bankRepository.Delete(bankID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Something went wrong"})
		return
	}

	c.JSON(http.StatusNoContent, nil)
}
