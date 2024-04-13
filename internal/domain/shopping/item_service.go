package shopping

import (
	"errors"
	"funbanking/internal/domain/users"
	"funbanking/package/utils"
	"strconv"
)

type ItemService interface {
	FindByID(itemID string) (Item, error)
	FindAllByShopID(shopID string) ([]Item, error)
	Save(item Item, user users.User) (Item, error)
	Delete(itemID string, user users.User) error
}

type itemService struct {
	itemRepository ItemRepository
	shopService    ShopService
}

func NewItemService(shopService ShopService) ItemService {
	return itemService{
		itemRepository: NewItemRepository(),
		shopService:    shopService,
	}
}

func (service itemService) FindByID(itemID string) (Item, error) {
	var item Item

	if err := service.itemRepository.FindByID(itemID, &item); err != nil {
		return Item{}, err
	}

	return item, nil
}

func (service itemService) FindAllByShopID(shopID string) ([]Item, error) {
	var items []Item

	if err := service.itemRepository.FindAllByShopID(shopID, &items); err != nil {
		return make([]Item, 0), err
	}

	return utils.Listify(items), nil
}

func (service itemService) Save(item Item, user users.User) (Item, error) {
	shop, err := service.shopService.FindByID(strconv.Itoa(int(item.ShopID)))

	if err != nil {
		return Item{}, err
	}

	if !service.hasAccess(shop, user) {
		return Item{}, errors.New("forbidden")
	}

	if len(item.Name) == 0 && len(item.Description) == 0 {
		return Item{}, errors.New("missing required fields")
	}

	if err := service.itemRepository.Save(&item); err != nil {
		return item, err
	}
	return item, nil
}

// todo this has bugs
func (service itemService) Delete(itemID string, user users.User) error {
	item, err := service.FindByID(itemID)

	if err != nil {
		return err
	}

	if item.Shop.UserID != user.ID {
		return errors.New("forbidden")
	}

	return service.itemRepository.Delete(itemID)
}

func (service itemService) hasAccess(shop Shop, user users.User) bool {
	return shop.UserID == user.ID
}
