package utils

type PaginatedResult[T any] struct {
	Data     []T   `json:"data"`
	Total    int64 `json:"total"`
	Page     int   `json:"page"`
	PageSize int   `json:"page_size"`
}

func NewPaginatedResult[T any](data []T, total int64, page, pageSize int) *PaginatedResult[T] {
	return &PaginatedResult[T]{
		Data:     data,
		Total:    total,
		Page:     page,
		PageSize: pageSize,
	}
}
