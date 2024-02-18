package utils

import (
	"funbanking/internal/domain/users"
	"funbanking/internal/infrastructure/persistence"
	"funbanking/package/constants"
)

func IsAdmin(userID string) bool {
	var user users.User

	if err := persistence.DB.Find(&user, "id = ?", userID).Error; err != nil {
		return false
	}

	return user.Role == constants.AdminRole
}
