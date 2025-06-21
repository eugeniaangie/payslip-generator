package model

const (
	DefaultPage     = 1
	DefaultPageSize = 10
	MaxPageSize     = 100
)

type Pagination struct {
	TotalPage   int64 `json:"total_page"`
	RowPerPage  int64 `json:"row_per_page"`
	TotalRow    int64 `json:"total_row"`
	StartRow    int64 `json:"start_row"`
	EndRow      int64 `json:"end_row"`
	CurrentPage int64 `json:"current_page"`
	HasNext     bool  `json:"has_next"`
}

type PaginationParam struct {
	Page     int64 `json:"page"`
	PageSize int64 `json:"page_size"`
}