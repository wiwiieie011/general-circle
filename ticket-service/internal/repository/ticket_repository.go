package repository

import (
	"ticket-service/internal/models"

	"gorm.io/gorm"
)

type TicketRepository struct {
	db *gorm.DB
}

func NewTicketRepository(db *gorm.DB) *TicketRepository {
	return &TicketRepository{db: db}
}

func (r *TicketRepository) Create(db *gorm.DB, ticket *models.Ticket) error {
	return db.Create(ticket).Error
}
