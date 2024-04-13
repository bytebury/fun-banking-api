package shopping

import (
	"errors"
	"funbanking/internal/domain/banking"
	"math"
)

type PurchaseService interface {
	BuyItems(items []Item, accountID string) ([]Purchase, error)
}

type purchaseService struct {
	shopService        ShopService
	purchaseRepository PurchaseRepository
	accountService     banking.AccountService
}

func NewPurchaseService() PurchaseService {
	return purchaseService{
		shopService: NewShopService(
			NewShopRepository(),
		),
		purchaseRepository: NewPurchaseRepository(),
		accountService: banking.NewAccountService(
			banking.NewAccountRepository(),
			banking.NewCustomerRepository(),
		),
	}
}

func (service purchaseService) BuyItems(items []Item, accountID string) ([]Purchase, error) {
	if len(items) == 0 {
		return make([]Purchase, 0), errors.New("empty cart")
	}

	account, accountErr := service.accountService.FindByID(accountID)
	if accountErr != nil {
		return make([]Purchase, 0), accountErr
	}

	cartPrice := service.calculateCartPrice(items)

	if account.Balance < cartPrice {
		return make([]Purchase, 0), errors.New("insufficient funds")
	}

	return service.purchaseRepository.BuyItems(items, math.Abs(cartPrice), account)
}

func (service purchaseService) calculateCartPrice(items []Item) float64 {
	totalPrice := float64(0)

	for _, item := range items {
		totalPrice += item.Price
	}

	return totalPrice
}
