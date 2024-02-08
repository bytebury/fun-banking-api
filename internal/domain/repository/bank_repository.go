package repository

import (
	"funbanking/internal/domain/model"
	"funbanking/internal/infrastructure/persistence"
	"strconv"

	"gorm.io/gorm"
)

type BankRepository interface {
	FindByID(bankID string, bank *model.Bank) error
	FindAllByUserID(userID string, banks *[]model.Bank) error
	FindByUsernameAndSlug(username, slug string, bank *model.Bank) error
	FindAllCustomers(bankID string, customers *[]model.Customer) error
	Create(bank *model.Bank) error
	Update(bank *model.Bank) error
	Delete(bankID string) error
}

type bankRepository struct {
	db *gorm.DB
}

func NewBankRepository() BankRepository {
	return bankRepository{db: persistence.DB}
}

func (r bankRepository) FindByID(bankID string, bank *model.Bank) error {
	return r.db.Preload("User").First(&bank, bankID).Error
}

func (r bankRepository) FindAllByUserID(userID string, banks *[]model.Bank) error {
	return r.db.Preload("User").Find(&banks, "user_id = ?", userID).Error
}

func (r bankRepository) FindByUsernameAndSlug(username, slug string, bank *model.Bank) error {
	var user model.User

	if err := r.db.First(&user, "username = ?", username).Error; err != nil {
		return err
	}

	return r.db.First(&bank, "user_id = ? AND slug = ?", user.ID, slug).Error
}

func (r bankRepository) FindAllCustomers(bankID string, customers *[]model.Customer) error {
	return r.db.Find(&customers, "bank_id = ?", bankID).Error
}

func (r bankRepository) Create(bank *model.Bank) error {
	return r.db.Create(&bank).Error
}

// TODO: This should be a transaction!
func (r bankRepository) Update(bank *model.Bank) error {
	var foundBank model.Bank

	if err := r.FindByID(strconv.Itoa(int(bank.ID)), &foundBank); err != nil {
		return err
	}

	if bank.Name == "" {
		bank.Name = foundBank.Name
	}

	if bank.Slug == "" {
		bank.Slug = foundBank.Slug
	}

	if bank.Description == "" {
		bank.Description = foundBank.Description
	}

	return r.db.Model(&bank).Select("Name", "Slug", "Description").Updates(&bank).Error
}

func (r bankRepository) Delete(bankID string) error {
	return r.db.Delete(&model.Bank{}, bankID).Error
}
