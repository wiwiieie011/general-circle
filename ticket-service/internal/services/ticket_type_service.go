package services

import (
	"context"
	api_http "ticket-service/internal/api/http"
	"ticket-service/internal/dto"
	dto_api "ticket-service/internal/dto/api"
	"ticket-service/internal/models"
	"ticket-service/internal/repository"
)

type TicketTypeService struct {
	eventClient    *api_http.EventClient
	ticketTypeRepo *repository.TicketTypeRepository
}

func NewTicketTypeService(
	eventClient *api_http.EventClient,
	ticketTypeRepo *repository.TicketTypeRepository,
) *TicketTypeService {
	return &TicketTypeService{
		eventClient:    eventClient,
		ticketTypeRepo: ticketTypeRepo,
	}
}

func (s *TicketTypeService) Create(
	ctx context.Context,
	eventId uint64,
	requestDto dto.CreateTicketTypeRequest,
) (*models.TicketType, error) {
	eventResp, err := s.eventClient.GetEvent(ctx, eventId)
	if err != nil {
		return nil, err
	}

	if eventResp.Status != dto_api.EventStatusPublished {
		return nil, dto.ErrEventNotPublished
	}

	ticketType := &models.TicketType{
		EventID:    eventId,
		Type:       models.TicketTypeKind(requestDto.Type),
		Price:      requestDto.Price,
		Quantity:   int(requestDto.Quantity),
		SalesStart: requestDto.SalesStart,
		SalesEnd:   requestDto.SalesEnd,
		Sold:       0,
	}

	if err := s.ticketTypeRepo.Create(ctx, ticketType); err != nil {
		return nil, err
	}

	return ticketType, nil
}
