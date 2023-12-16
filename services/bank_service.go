package services

import (
	"golfer/models"
	"golfer/repositories"

	"github.com/gin-gonic/gin"
)

type BankService struct {
	bankRepository repositories.BankRepository
}

func NewBankService(bankRepository repositories.BankRepository) *BankService {
	return &BankService{
		bankRepository,
	}
}

func (service BankService) Create(request *models.Bank) error {
	return service.bankRepository.Create(request)
}

func (service BankService) FindByID(bankID string, bank *models.Bank) error {
	return service.bankRepository.FindByID(bankID, bank)
}

func (service BankService) Search(c *gin.Context, banks *[]models.Bank) error {
	return service.bankRepository.Search(c, banks)
}

func (service BankService) Update(bankID string, request *models.Bank) (models.Bank, error) {
	var bank models.Bank

	if err := service.bankRepository.FindByID(bankID, &bank); err != nil {
		return bank, err
	}

	if request.Slug != "" {
		bank.Slug = request.Slug
	}

	if request.Name != "" {
		bank.Name = request.Name
	}

	if request.Description != "" {
		bank.Description = request.Description
	}

	if err := service.bankRepository.Update(&bank); err != nil {
		return bank, err
	}

	return bank, nil
}

func (service BankService) Delete(bankID string) error {
	return service.bankRepository.Delete(bankID)
}
