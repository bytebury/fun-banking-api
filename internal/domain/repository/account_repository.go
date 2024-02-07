package repository

import (
	"funbanking/internal/infrastructure/persistence"

	"gorm.io/gorm"
)

type AccountRepository interface {

}

type accountRepository struct {
	db *gorm.DB
}

func NewAccountRepository() AccountRepository {
	return accountRepository{ db: persistence.DB }
}
