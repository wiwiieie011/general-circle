package dto

type Status string

const (
	Draft     Status = "draft"
	Published Status = "published"
	Cancelled Status = "cancelled"
)

type CreateEventRequest struct {
	Title      string `json:"title" binding:"required"`
	Seats      *int   `json:"seats"`
	UserID     uint   `json:"user_id" binding:"required"`
	CategoryID *uint  `json:"category_id"`
}

type UpdateEventRequest struct {
	Title      *string `json:"title"`
	Seats      *int    `json:"seats"`
	UserID     *uint   `json:"user_id"`
	CategoryID *uint   `json:"category_id"`
}
