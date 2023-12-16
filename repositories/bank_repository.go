package repositories

import (
	"golfer/database"
	"golfer/models"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type BankRepository struct {
	db *gorm.DB
}

func NewBankRepository() *BankRepository {
	return &BankRepository{
		db: database.DB,
	}
}

func (repository BankRepository) Create(bank *models.Bank) error {
	return repository.db.Create(&bank).Error
}

func (repository BankRepository) FindByID(bankID string, bank *models.Bank) error {
	return repository.db.Preload("User").First(&bank, bankID).Error
}

func (repository BankRepository) Search(c *gin.Context, banks *[]models.Bank) error {
	id := c.Query("id")
	name := c.Query("name")
	slug := c.Query("slug")
	// TODO(Marcello): You should only be able to search for your own banks
	ownerID := c.Query("owner-id")
	limit, err := strconv.Atoi(c.Query("limit"))

	query := repository.db

	if id == "" && name == "" && slug == "" && ownerID == "" {
		*banks = make([]models.Bank, 0)
		return nil
	}

	if id != "" {
		query = query.Where("id = ?", id)
	}

	if name != "" {
		query = query.Where("name ILIKE ?", "%"+name+"%")
	}

	if slug != "" {
		query = query.Where("slug = ?", slug)
	}

	if ownerID != "" {
		query = query.Where("user_id = ?", ownerID)
	}

	if limit == 0 || err != nil {
		limit = 25
	}

	return query.Preload("User").Limit(limit).Find(&banks).Error
}

func (repository BankRepository) Update(bank *models.Bank) error {
	return repository.db.Save(&bank).Error
}

func (repository BankRepository) Delete(bankID string) error {
	return repository.db.Delete(&models.Bank{}, "id = ?", bankID).Error
}
