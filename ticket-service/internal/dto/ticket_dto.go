package dto

type CreateTicketRequest struct {
	TicketTypeID uint64 `json:"ticket_type_id" binding:"required,gt=0"`
	UserID       uint64 `json:"user_id" binding:"required,gt=0"`
}
