package models

type PagingInfo struct {
	PageNumber uint `json:"page_number"`
	TotalItems uint `json:"total_items"`
}

type PaginatedResponse[T any] struct {
	Items      []T        `json:"items"`
	PagingInfo PagingInfo `json:"paging_info"`
}
