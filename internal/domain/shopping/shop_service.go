package shopping

import (
	"errors"
	"funbanking/internal/domain/users"
	"funbanking/package/constants"
	"funbanking/package/utils"
)

type ShopService interface {
	FindByID(shopID string) (Shop, error)
	FindAllByUser(user users.User) ([]Shop, error)
	Save(shop Shop, user users.User) (Shop, error)
	Delete(shopID string, user users.User) error
}

type shopService struct {
	shopRepository ShopRepository
}

func NewShopService(shopRepository ShopRepository) ShopService {
	return shopService{shopRepository}
}

func (service shopService) FindByID(shopID string) (Shop, error) {
	var shop Shop

	if err := service.shopRepository.FindByID(shopID, &shop); err != nil {
		return Shop{}, err
	}

	return shop, nil
}

func (service shopService) FindAllByUser(user users.User) ([]Shop, error) {
	var shops []Shop

	if err := service.shopRepository.FindAllByUser(user, &shops); err != nil {
		return make([]Shop, 0), err
	}

	return utils.Listify[Shop](shops), nil
}

func (service shopService) Save(shop Shop, user users.User) (Shop, error) {
	if shop.ID == 0 {
		shop.UserID = user.ID
	}

	if !service.hasAccess(shop, user) {
		return shop, errors.New("forbidden")
	}

	if err := service.shopRepository.Save(&shop); err != nil {
		return shop, err
	}
	return shop, nil
}

func (service shopService) Delete(shopID string, user users.User) error {
	var shop Shop
	if err := service.shopRepository.FindByID(shopID, &shop); err != nil {
		return err
	}

	if !service.hasAccess(shop, user) {
		return errors.New("forbidden")
	}

	return service.shopRepository.Delete(shopID)
}

func (service shopService) hasAccess(shop Shop, user users.User) bool {
	return shop.UserID == user.ID || user.Role == constants.AdminRole
}
