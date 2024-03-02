package banking

import (
	"errors"
	"funbanking/internal/infrastructure/persistence"

	"gorm.io/gorm"
)

type CustomerRepository interface {
	FindByID(customerID string, customer *Customer) error
	FindAccounts(customerID string, accounts *[]Account) error
	FindByBankAndPIN(bankID string, pin string, customer *Customer) error
	Create(customer *Customer) error
	Update(customerID string, customer *Customer) error
	Delete(customerID string) error
}

type customerRepository struct {
	db *gorm.DB
}

func NewCustomerRepository() CustomerRepository {
	return customerRepository{db: persistence.DB}
}

func (r customerRepository) FindByID(customerID string, customer *Customer) error {
	return r.db.Preload("Accounts").Preload("Bank").First(&customer, "id = ?", customerID).Error
}

func (r customerRepository) FindAccounts(customerID string, accounts *[]Account) error {
	return r.db.Find(&accounts, "customer_id = ?", customerID).Error
}

func (r customerRepository) FindByBankAndPIN(bankID string, pin string, customer *Customer) error {
	if bankID == "" || pin == "" {
		return errors.New("missing required fields: BankID or PIN")
	}

	return r.db.Preload("Accounts").First(&customer, "bank_id = ? AND pin = ?", bankID, pin).Error
}

func (r customerRepository) Create(customer *Customer) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&customer).Error; err != nil {
			return err
		}

		if err := tx.Create(&Account{Name: "Checkings", CustomerID: customer.ID}).Error; err != nil {
			return err
		}

		return nil
	})
}

func (r customerRepository) Update(customerID string, customer *Customer) error {
	var foundCustomer Customer

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
	return r.db.Delete(&Customer{}, customerID).Error
}
