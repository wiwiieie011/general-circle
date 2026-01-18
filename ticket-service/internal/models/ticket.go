package models

type TicketStatus string

const (
	TicketStatusActive    TicketStatus = "active"
	TicketStatusUsed      TicketStatus = "used"
	TicketStatusCancelled TicketStatus = "cancelled"
	TicketStatusExpired   TicketStatus = "expired"
)

type Ticket struct {
	Base
	EventID      uint64       `json:"event_id" gorm:"not null"`
	TicketTypeID uint         `json:"ticket_type_id" gorm:"not null"`
	TicketType   TicketType   `json:"-" gorm:"foreignKey:TicketTypeID;references:ID"`
	UserID       uint64       `json:"user_id" gorm:"not null"`
	Code         string       `json:"code" gorm:"type:varchar(64);not null;uniqueIndex"`
	Status       TicketStatus `json:"status" gorm:"type:varchar(20);not null;default:'active'"`
}
