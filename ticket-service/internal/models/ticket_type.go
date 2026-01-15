package models

import "time"

type TicketTypeKind string

const (
	TicketTypeStandard  TicketTypeKind = "standard"
	TicketTypeVIP       TicketTypeKind = "vip"
	TicketTypeEarlyBird TicketTypeKind = "early_bird"
)

type TicketType struct {
	Base
	EventID    uint64         `json:"event_id" gorm:"not null"`
	Type       TicketTypeKind `json:"type" gorm:"type:varchar(20);not null"`
	Price      int64          `json:"price" gorm:"not null"`
	Quantity   int            `json:"quantity" gorm:"not null"`
	Sold       int            `json:"sold" gorm:"not null;default:0"`
	SalesStart time.Time      `json:"sales_start" gorm:"not null"`
	SalesEnd   time.Time      `json:"sales_end" gorm:"not null"`
}
