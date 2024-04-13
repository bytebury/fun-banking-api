package shopping

import (
	"funbanking/internal/infrastructure/persistence"

	"gorm.io/gorm"
)

type ItemRepository interface {
	FindByID(itemID string, item *Item) error
	FindAllByShopID(shopID string, items *[]Item) error
	Save(item *Item) error
	Delete(itemID string) error
}

type itemRepository struct {
	db *gorm.DB
}

func NewItemRepository() ItemRepository {
	return itemRepository{
		db: persistence.DB,
	}
}

func (repo itemRepository) FindByID(itemID string, item *Item) error {
	return repo.db.First(&item, "id = ?", itemID).Error
}

func (repo itemRepository) FindAllByShopID(shopID string, items *[]Item) error {
	return repo.db.Find(&items, "shop_id = ?", shopID).Error
}

func (repo itemRepository) Save(item *Item) error {
	return repo.db.Save(&item).Error
}

func (repo itemRepository) Delete(itemID string) error {
	return repo.db.Delete(&Item{}, "id = ?", itemID).Error
}
