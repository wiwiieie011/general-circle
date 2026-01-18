package repository

import (
	"context"
	"ticket-service/internal/models"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
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

func (r *TicketTypeRepository) GetByID(ctx context.Context, id uint64) (*models.TicketType, error) {
	var ticketType models.TicketType

	err := r.db.WithContext(ctx).First(&ticketType, id).Error

	if err != nil {
		return nil, err
	}

	return &ticketType, nil
}

func (r *TicketTypeRepository) GetByIDForUpdate(db *gorm.DB, id uint64) (*models.TicketType, error) {
	var tt models.TicketType
	err := db.
		Clauses(clause.Locking{Strength: "UPDATE"}).
		First(&tt, id).
		Error
	return &tt, err
}

func (r *TicketTypeRepository) IncrementSold(db *gorm.DB, id uint64) error {
	return db.Model(&models.TicketType{}).
		Where("id = ?", id).
		UpdateColumn("sold", gorm.Expr("sold + 1")).
		Error
}
