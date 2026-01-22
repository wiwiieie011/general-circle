package repository

import (
	"errors"
	"ticket-service/internal/dto"
	"ticket-service/internal/models"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type TicketRepository struct {
	db *gorm.DB
}

func NewTicketRepository(db *gorm.DB) *TicketRepository {
	return &TicketRepository{db: db}
}

func (r *TicketRepository) WithDB(db *gorm.DB) *TicketRepository {
	return &TicketRepository{db: db}
}

func (r *TicketRepository) Create(ticket *models.Ticket) error {
	return r.db.Create(ticket).Error
}

func (r *TicketRepository) List(filter dto.TicketListFilter) ([]models.Ticket, error) {
	var tickets []models.Ticket

	q := r.db.Model(&models.Ticket{})

	if filter.EventID != nil {
		q = q.Where("event_id = ?", *filter.EventID)
	}

	if filter.UserID != nil {
		q = q.Where("user_id = ?", *filter.UserID)
	}

	if filter.Status != nil {
		q = q.Where("status = ?", *filter.Status)
	}

	if filter.Limit != nil && *filter.Limit > 0 {
		q = q.Limit(*filter.Limit)
	}

	if filter.Offset != nil && *filter.Offset > 0 {
		q = q.Offset(*filter.Offset)
	}

	if err := q.Order("created_at DESC").Find(&tickets).Error; err != nil {
		return nil, err
	}

	return tickets, nil
}

func (r *TicketRepository) IsExist(code string) (bool, error) {
	var ticket models.Ticket

	err := r.db.
		Model(&models.Ticket{}).
		Where("code = ?", code).
		Where("status = ?", models.TicketStatusActive).
		First(&ticket).
		Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, nil
		}
		return false, err
	}

	return true, nil
}

func (r *TicketRepository) Checkin(code string) (*models.Ticket, error) {
	var ticket models.Ticket

	err := r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.
			Clauses(clause.Locking{Strength: "UPDATE"}).
			Where("code = ?", code).
			Where("status = ?", models.TicketStatusActive).
			First(&ticket).
			Error; err != nil {

			if errors.Is(err, gorm.ErrRecordNotFound) {
				return dto.ErrTicketNotFoundOrNotActive
			}
			return err
		}

		if err := tx.
			Model(&ticket).
			Update("status", models.TicketStatusUsed).
			Error; err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return &ticket, nil
}
