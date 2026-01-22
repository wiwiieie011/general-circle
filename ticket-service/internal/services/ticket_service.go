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

		ticketRepo := s.ticketRepo.WithDB(tx)
		ticketTypeRepo := s.ticketTypeRepo.WithDB(tx)

		ticketType, err := ticketTypeRepo.GetByIDForUpdate(requestDto.TicketTypeID)
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

		if err := ticketTypeRepo.IncrementSold(requestDto.TicketTypeID); err != nil {
			return err
		}

		ticket = &models.Ticket{
			EventID:      eventId,
			TicketTypeID: ticketType.ID,
			UserID:       requestDto.UserID,
			Code:         uuid.NewString(),
			Status:       models.TicketStatusActive,
		}

		return ticketRepo.Create(ticket)
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
		s.logger.Warn("kafka publish failed", "error", err.Error())
	}

	return ticket, nil
}

func (s *TicketService) List(filter dto.TicketListFilter) ([]models.Ticket, error) {
	tickets, err := s.ticketRepo.List(filter)
	if err != nil {
		s.logger.Error(err.Error())
		return tickets, err
	}

	return tickets, nil
}

func (s *TicketService) IsExist(codeDto *dto.TicketCode) (bool, error) {
	isExist, err := s.ticketRepo.IsExist(codeDto.Code)
	if err != nil {
		return false, err
	}

	return isExist, nil
}

func (s *TicketService) Checkin(ctx context.Context, codeDto *dto.TicketCode) error {
	ticket, err := s.ticketRepo.Checkin(codeDto.Code)
	if err != nil {
		return err
	}

	event := kafka_events.TicketCheckinEvent{
		TicketID:     uint64(ticket.ID),
		EventID:      ticket.EventID,
		TicketTypeID: uint64(ticket.TicketTypeID),
		UserID:       ticket.UserID,
		CheckedinAt:  time.Now(),
	}

	if err := s.kafkaProducer.PublishTicketCheckin(ctx, event); err != nil {
		s.logger.Warn("kafka publish failed", "error", err.Error())
	}

	return nil
}
