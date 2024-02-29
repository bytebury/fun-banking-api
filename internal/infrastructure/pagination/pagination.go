package pagination

import (
	"gorm.io/gorm"
)

type PaginatedResponse[T any] struct {
	Items        []T `json:"items"`
	PageNumber   int `json:"page_number"`
	ItemsPerPage int `json:"items_per_page"`
	TotalItems   int `json:"total_items"`
}

/**
 * This will take in a query that is already ordered and filtered and convert it to a
 * Paginated response for ease of use.
 */
func Find[T any](query *gorm.DB, pageNumber int, itemsPerPage int) (PaginatedResponse[T], error) {
	if pageNumber == 0 {
		pageNumber = 1
	}

	if itemsPerPage == 0 {
		itemsPerPage = 25
	}

	var results []T
	var count int64

	offset := (pageNumber - 1) * itemsPerPage

	if err := query.Model(&results).Count(&count).Error; err != nil {
		return PaginatedResponse[T]{
			Items:        make([]T, 0),
			ItemsPerPage: itemsPerPage,
			PageNumber:   pageNumber,
			TotalItems:   int(count),
		}, err
	}

	if err := query.Limit(itemsPerPage).Offset(offset).Find(&results).Error; err != nil {
		return PaginatedResponse[T]{
			Items:        make([]T, 0),
			ItemsPerPage: itemsPerPage,
			PageNumber:   pageNumber,
			TotalItems:   int(count),
		}, err
	}

	return PaginatedResponse[T]{
		Items:        results,
		ItemsPerPage: itemsPerPage,
		PageNumber:   pageNumber,
		TotalItems:   int(count),
	}, nil
}
