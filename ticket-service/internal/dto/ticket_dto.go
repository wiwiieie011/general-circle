package dto

import "ticket-service/internal/models"

type CreateTicketRequest struct {
	TicketTypeID uint64 `json:"ticket_type_id" binding:"required,gt=0"`
	UserID       uint64 `json:"user_id" binding:"required,gt=0"`
}

type TicketListFilter struct {
	EventID *uint64 `json:"event_id"`
	UserID  *uint64 `json:"user_id"`

	Status *models.TicketStatus `json:"status"`

	Limit  *int `json:"limit"`
	Offset *int `json:"offset"`
}

type TicketCode struct {
	Code string `json:"code"`
}
