package repository

import (
	"context"
	"ticket-service/internal/models"

	"gorm.io/gorm"
)

type TicketTypeRepository struct {
	db *gorm.DB
}

func NewTicketTypeRepository(db *gorm.DB) *TicketTypeRepository {
	return &TicketTypeRepository{db: db}
}

func (r *TicketTypeRepository) Create(ctx context.Context, ticketType *models.TicketType) error {
	return r.db.WithContext(ctx).Create(ticketType).Error
}
