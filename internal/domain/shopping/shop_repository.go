package shopping

import (
	"funbanking/internal/domain/users"
	"funbanking/internal/infrastructure/persistence"

	"gorm.io/gorm"
)

type ShopRepository interface {
	FindByID(shopID string, shop *Shop) error
	FindAllByUser(user users.User, shops *[]Shop) error
	Save(shop *Shop) error
	Delete(shopID string) error
}

type shopRepository struct {
	db *gorm.DB
}

func NewShopRepository() ShopRepository {
	return shopRepository{db: persistence.DB}
}

func (repo shopRepository) FindByID(shopID string, shop *Shop) error {
	return repo.db.First(&shop, "id = ?", shopID).Error
}

func (repo shopRepository) FindAllByUser(user users.User, shops *[]Shop) error {
	return repo.db.Find(&shops, "user_id = ?", user.ID).Error
}

func (repo shopRepository) Save(shop *Shop) error {
	return repo.db.Save(&shop).Error
}

func (repo shopRepository) Delete(shopID string) error {
	return repo.db.Delete(&Shop{}, "id = ?", shopID).Error
}
