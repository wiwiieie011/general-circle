package kafka

import "time"

type TicketPurchasedEvent struct {
	TicketID     uint64    `json:"ticket_id"`
	EventID      uint64    `json:"event_id"`
	TicketTypeID uint64    `json:"ticket_type_id"`
	UserID       uint64    `json:"user_id"`
	Code         string    `json:"code"`
	Status       string    `json:"status"`
	PurchasedAt  time.Time `json:"purchased_at"`
}

type TicketCheckinEvent struct {
	TicketID     uint64    `json:"ticket_id"`
	EventID      uint64    `json:"event_id"`
	TicketTypeID uint64    `json:"ticket_type_id"`
	UserID       uint64    `json:"user_id"`
	CheckedinAt  time.Time `json:"checked_in_at"`
}
