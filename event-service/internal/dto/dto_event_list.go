package dto

const (
	DefaultPage      = 1
	DefaultLimit     = 10
)

type EventListQuery struct {
	// Фильтры
	Title  string `form:"title"`
	Status string `form:"status"`

	// Пагинация
	Page  int `form:"page"`
	Limit int `form:"limit"`

	// Сортивка
	// sort_by: created_at и title
	// sort_order: asc и desc
	SortBy    string `form:"sort_by"`
	SortOrder string `form:"sort_order"`
}
