package banking

import (
	"funbanking/internal/domain/users"
	"funbanking/internal/infrastructure/persistence"
	"regexp"
	"strings"

	"gorm.io/gorm"
)

type BankRepository interface {
	FindByID(id string, bank *Bank) error
	FindByUsernameAndSlug(username, slug string, bank *Bank) error
	FindAllCustomers(bankID string, customers *[]Customer) error
	FindAllByUserID(userID string, banks *[]Bank) error
	Create(bank *Bank) error
	Update(bankID string, bank *Bank) error
	Delete(bankID string) error
}

type bankRepository struct {
	db *gorm.DB
}

func NewBankRepository() BankRepository {
	return bankRepository{db: persistence.DB}
}

func (r bankRepository) FindByID(bankID string, bank *Bank) error {
	return r.db.Preload("User").First(&bank, bankID).Error
}

func (r bankRepository) FindAllByUserID(userID string, banks *[]Bank) error {
	return r.db.Preload("User").Find(&banks, "user_id = ?", userID).Error
}

func (r bankRepository) FindByUsernameAndSlug(username, slug string, bank *Bank) error {
	var user users.User

	if err := r.db.First(&user, "username = ?", username).Error; err != nil {
		return err
	}

	return r.db.First(&bank, "user_id = ? AND slug = ?", user.ID, slug).Error
}

func (r bankRepository) FindAllCustomers(bankID string, customers *[]Customer) error {
	return r.db.Preload("Accounts").Order("last_name ASC, first_name ASC").Find(&customers, "bank_id = ?", bankID).Error
}

func (r bankRepository) Create(bank *Bank) error {
	r.normalize(bank)
	return r.db.Create(&bank).Error
}

func (r bankRepository) Update(bankID string, bank *Bank) error {
	var foundBank Bank

	if err := r.FindByID(bankID, &foundBank); err != nil {
		return err
	}

	if bank.Name == "" {
		bank.Name = foundBank.Name
	}

	if bank.Description == "" {
		bank.Description = foundBank.Description
	}

	r.normalize(bank)

	return r.db.Model(&foundBank).Select("Name", "Slug", "Description").Updates(&bank).Error
}

func (r bankRepository) Delete(bankID string) error {
	return r.db.Delete(&Bank{}, bankID).Error
}

func (r bankRepository) normalize(bank *Bank) {
	bank.Name = strings.TrimSpace(bank.Name)
	bank.Slug = r.getSlugFromName(strings.ToLower(bank.Name))
}

func (r bankRepository) getSlugFromName(name string) string {
	re := regexp.MustCompile(`\s+`)
	return re.ReplaceAllString(name, "-")
}
