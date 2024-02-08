package repository

import (
	"funbanking/internal/domain/model"
	"funbanking/internal/infrastructure/persistence"

	"gorm.io/gorm"
)

type CustomerRepository interface {
	FindByID(customerID string, customer *model.Customer) error
	FindAccounts(customerID string, accounts *[]model.Account) error
	Create(customer *model.Customer) error
	Update(customerID string, customer *model.Customer) error
	Delete(customerID string) error
}

type customerRepository struct {
	db *gorm.DB
}

func NewCustomerRepository() CustomerRepository {
	return customerRepository{db: persistence.DB}
}

func (r customerRepository) FindByID(customerID string, customer *model.Customer) error {
	return r.db.Preload("Accounts").First(&customer, "id = ?", customerID).Error
}

func (r customerRepository) FindAccounts(customerID string, accounts *[]model.Account) error {
	return r.db.Find(&accounts, "customer_id = ?", customerID).Error
}

func (r customerRepository) Create(customer *model.Customer) error {
	// When you create a customer, you also create a checkings account
	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&customer).Error; err != nil {
			return err
		}

		if err := tx.Create(&model.Account{Name: "Checkings", CustomerID: customer.ID}).Error; err != nil {
			return err
		}

		return nil
	})
}

func (r customerRepository) Update(customerID string, customer *model.Customer) error {
	var foundCustomer model.Customer

	if err := r.FindByID(customerID, &foundCustomer); err != nil {
		return err
	}

	if customer.FirstName == "" {
		customer.FirstName = foundCustomer.FirstName
	}

	if customer.LastName == "" {
		customer.LastName = foundCustomer.LastName
	}

	if customer.PIN == "" {
		customer.PIN = foundCustomer.PIN
	}

	return r.db.Model(&foundCustomer).Select("FirstName", "LastName", "PIN").Updates(&customer).Error
}

func (r customerRepository) Delete(customerID string) error {
	return r.db.Delete(&model.Customer{}, customerID).Error
}
