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
	ID           uint64       `gorm:"primaryKey"`
	EventID      uint64       `gorm:"not null"`
	TicketTypeID uint64       `gorm:"not null"`
	UserID       uint64       `gorm:"not null"`
	Code         string       `gorm:"type:varchar(64);not null;uniqueIndex"`
	Status       TicketStatus `gorm:"type:varchar(20);not null"`
}
