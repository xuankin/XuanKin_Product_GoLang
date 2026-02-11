package models

type APIResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
	Errors  interface{} `json:"errors,omitempty"`
}

type PaginationResponse struct {
	Total       int64       `json:"total"`
	CurrentPage int         `json:"current_page"`
	LastPage    int         `json:"last_page"`
	Data        interface{} `json:"data"`
}
type FilterParams struct {
	Page      int    `form:"page,default=1"`
	Limit     int    `form:"limit,default=10"`
	Search    string `form:"search"`
	SortBy    string `form:"sort_by"`
	SortOrder string `form:"sort_order"`
}
