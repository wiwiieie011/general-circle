package services

import (
	"context"
	"log/slog"
	api_http "ticket-service/internal/api/http"
	"ticket-service/internal/kafka"
	kafka_events "ticket-service/internal/kafka/events"

	"ticket-service/internal/dto"

	dto_api "ticket-service/internal/dto/api"
	"ticket-service/internal/models"
	"ticket-service/internal/repository"

	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type TicketService struct {
	ticketRepo     *repository.TicketRepository
	ticketTypeRepo *repository.TicketTypeRepository
	eventClient    *api_http.EventClient
	kafkaProducer  *kafka.Producer
	db             *gorm.DB
	logger         *slog.Logger
}

func NewTicketService(
	ticketRepo *repository.TicketRepository,
	ticketTypeRepo *repository.TicketTypeRepository,
	eventClient *api_http.EventClient,
	kafkaProducer *kafka.Producer,
	db *gorm.DB,
	logger *slog.Logger,
) *TicketService {
	return &TicketService{
		ticketRepo:     ticketRepo,
		ticketTypeRepo: ticketTypeRepo,
		eventClient:    eventClient,
		kafkaProducer:  kafkaProducer,
		db:             db,
		logger:         logger,
	}
}

func (s *TicketService) Create(ctx context.Context, eventId uint64, requestDto dto.CreateTicketRequest) (*models.Ticket, error) {
	eventResp, err := s.eventClient.GetEvent(ctx, eventId)
	if err != nil {
		return nil, err
	}

	if eventResp.Status != dto_api.EventStatusPublished {
		return nil, dto.ErrEventNotPublished
	}

	var ticket *models.Ticket

	err = s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {

		ticketType, err := s.ticketTypeRepo.GetByIDForUpdate(tx, requestDto.TicketTypeID)
		if err != nil {
			return err
		}

		now := time.Now()

		if now.Before(ticketType.SalesStart) {
			return dto.ErrEventNotStarted
		}

		if now.After(ticketType.SalesEnd) {
			return dto.ErrEventEnded
		}

		if ticketType.Sold >= ticketType.Quantity {
			return dto.ErrTicketSoldOut
		}

		if err := s.ticketTypeRepo.IncrementSold(tx, requestDto.TicketTypeID); err != nil {
			return err
		}

		ticket = &models.Ticket{
			EventID:      eventId,
			TicketTypeID: ticketType.ID,
			UserID:       requestDto.UserID,
			Code:         uuid.NewString(),
			Status:       models.TicketStatusActive,
		}

		err = s.ticketRepo.Create(tx, ticket)

		return nil
	})

	if err != nil {
		s.logger.Error(err.Error())
		return nil, err
	}

	event := kafka_events.TicketPurchasedEvent{
		TicketID:     uint64(ticket.ID),
		EventID:      eventId,
		TicketTypeID: uint64(ticket.TicketTypeID),
		UserID:       ticket.UserID,
		Code:         ticket.Code,
		Status:       string(ticket.Status),
		PurchasedAt:  ticket.CreatedAt,
	}

	if err := s.kafkaProducer.PublishTicketPurchased(ctx, event); err != nil {
		s.logger.Warn("kafka publish failed", err.Error())
	}

	return ticket, nil
}
