package shopping

import (
	"errors"
	"fmt"
	"funbanking/internal/domain/banking"
	"funbanking/internal/infrastructure/persistence"
	"funbanking/package/utils"
	"strconv"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type PurchaseRepository interface {
	BuyItems(items []Item, totalPriceWithTax float64, account banking.Account) ([]Purchase, error)
}

type purchaseRepository struct {
	db                 *gorm.DB
	transactionService banking.TransactionService
	itemRepository     ItemRepository
}

func NewPurchaseRepository() PurchaseRepository {
	return purchaseRepository{
		db: persistence.DB,
		transactionService: banking.NewTransactionService(
			banking.NewTransactionRepository(),
		),
		itemRepository: NewItemRepository(),
	}
}

func (repo purchaseRepository) BuyItems(items []Item, cartPrice float64, account banking.Account) ([]Purchase, error) {
	cartID := uuid.New().String()

	// Buy each item individually
	err := repo.db.Transaction(func(tx *gorm.DB) error {
		var shop Shop

		// Find the shop
		if err := repo.db.First(&shop, "id = ?", items[0].ShopID).Error; err != nil {
			return err
		}

		totalTax := shop.TaxRate * cartPrice
		totalPrice := cartPrice + totalTax

		for _, item := range items {
			if item.NumberInStock == 0 {
				return errors.New("out of stock")
			}

			item.NumberInStock -= 1

			if err := repo.itemRepository.Save(&item); err != nil {
				return err
			}

			purchase := Purchase{
				CartID:     cartID,
				Item:       item,
				ItemID:     item.ID,
				Price:      item.Price,
				CartPrice:  cartPrice,
				TaxRate:    shop.TaxRate,
				TotalTax:   totalTax,
				TotalPrice: totalPrice,
			}
			if err := repo.db.Save(&purchase).Error; err != nil {
				return err
			}
		}

		// Take out the required funds from the customer via transactions
		transaction := banking.Transaction{
			Description:    fmt.Sprintf("Bought %d item(s) at %s", len(items), shop.Name),
			Amount:         totalPrice * -1,
			AccountID:      account.ID,
			CurrentBalance: account.Balance,
			Type:           "shopping",
			Origin:         banking.TransactionShopping,
		}

		if err := repo.transactionService.Create("", &transaction); err != nil {
			return err
		}

		if _, err := repo.transactionService.Approve("", strconv.Itoa(int(transaction.ID))); err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return make([]Purchase, 0), err
	}

	var purchases []Purchase
	if err := repo.db.Find(&purchases, "cart_id = ?", cartID).Error; err != nil {
		return make([]Purchase, 0), err
	}

	return utils.Listify(purchases), nil
}
