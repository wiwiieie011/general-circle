package dto

import "time"

type CreateTicketTypeRequest struct {
	Type       string    `json:"type" binding:"required,oneof=standard vip early_bird"`
	Price      int64     `json:"price" binding:"required,gt=0"`
	Quantity   int64     `json:"quantity" binding:"required,gt=0"`
	SalesStart time.Time `json:"sales_start" binding:"required"`
	SalesEnd   time.Time `json:"sales_end" binding:"required"`
}
