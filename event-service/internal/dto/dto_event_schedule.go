package dto

type Status string

const (
	Draft     Status = "draft"
	Published Status = "published"
)

type CreateEventRequest struct {
	Title  string `json:"title"`
	Status Status `json:"status"`
	Seats  *int   `json:"seats"`
	UserID uint   `json:"user_id"`
}

type UpdateEventRequest struct {
	Title  *string `json:"title"`
	Status *Status `json:"status"`
	Seats  *int    `json:"seats"`
	UserID *uint   `json:"user_id"`
}
