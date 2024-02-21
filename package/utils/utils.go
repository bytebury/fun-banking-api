package utils

import "strconv"

func StringToUIntPointer(value string) (*uint, error) {
	userIDAsInt, err := strconv.Atoi(value)

	if err != nil {
		return nil, err
	}

	userIDAsUInt := uint(userIDAsInt)
	return &userIDAsUInt, nil
}
