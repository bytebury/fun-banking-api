package utils

import (
	"funbanking/internal/domain/model"
	"funbanking/internal/infrastructure/persistence"
	"funbanking/package/constants"
)

func IsAdmin(userID string) bool {
	var user model.User

	if err := persistence.DB.Find(&user, "id = ?", userID).Error; err != nil {
		return false
	}

	return user.Role == constants.AdminRole
}
