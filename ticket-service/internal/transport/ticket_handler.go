package transport

import (
	"errors"
	"log/slog"
	"net/http"
	"strconv"
	"ticket-service/internal/dto"
	"ticket-service/internal/services"

	"github.com/gin-gonic/gin"
)

type TicketHandler struct {
	ticketTypeService *services.TicketTypeService
	ticketService     *services.TicketService
	logger            *slog.Logger
}

func NewTicketHandler(
	ticketTypeService *services.TicketTypeService,
	ticketService *services.TicketService,
	logger *slog.Logger,
) *TicketHandler {
	return &TicketHandler{
		ticketTypeService: ticketTypeService,
		ticketService:     ticketService,
		logger:            logger,
	}
}

func (h *TicketHandler) RegisterRoutes(r *gin.Engine) {
	r.GET("/ping", h.Ping)
	r.POST("/tickets/validate", h.TicketValidate)
	r.POST("/tickets/checkin", h.TicketCheckin)
	r.GET("/tickets", h.GetTickets)
	r.POST("/events/:id/ticket-types", h.CreateTicketType)
	r.POST("/events/:id/tickets", h.CreateTicket)
}

func (h *TicketHandler) Ping(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "success"})
}

func (h *TicketHandler) CreateTicketType(c *gin.Context) {
	ctx := c.Request.Context()
	eventId, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil || eventId <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid event id"})
		return
	}
	var ttDto dto.CreateTicketTypeRequest
	if err := c.ShouldBindJSON(&ttDto); err != nil {
		h.logger.Error(err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	ticketType, err := h.ticketTypeService.Create(ctx, uint64(eventId), ttDto)
	if err != nil {
		h.logger.Error(err.Error())
		switch {
		case errors.Is(err, dto.ErrEventNotFound):
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		case errors.Is(err, dto.ErrEventNotPublished):
			c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": ticketType})
}

func (h *TicketHandler) CreateTicket(c *gin.Context) {
	ctx := c.Request.Context()
	eventId, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil || eventId <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid event id"})
		return
	}
	var requestDto dto.CreateTicketRequest
	if err := c.ShouldBindJSON(&requestDto); err != nil {
		h.logger.Error(err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ticket, err := h.ticketService.Create(ctx, eventId, requestDto)

	if err != nil {
		h.logger.Error(err.Error())
		switch {
		case errors.Is(err, dto.ErrEventNotFound):
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		case errors.Is(err, dto.ErrEventNotPublished),
			errors.Is(err, dto.ErrTicketSoldOut),
			errors.Is(err, dto.ErrEventNotStarted),
			errors.Is(err, dto.ErrEventEnded):
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": ticket})
}

func (h *TicketHandler) GetTickets(c *gin.Context) {
	var filter dto.TicketListFilter
	if err := c.ShouldBindQuery(&filter); err != nil {
		h.logger.Error("binding error", "error", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	tickets, err := h.ticketService.List(filter)
	if err != nil {
		h.logger.Error(err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": tickets})
}

func (h *TicketHandler) TicketValidate(c *gin.Context) {
	var codeDto *dto.TicketCode
	if err := c.ShouldBindJSON(&codeDto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	isExist, err := h.ticketService.IsExist(codeDto)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"valid": isExist})
}

func (h *TicketHandler) TicketCheckin(c *gin.Context) {
	var codeDto *dto.TicketCode
	if err := c.ShouldBindJSON(&codeDto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := h.ticketService.Checkin(c.Request.Context(), codeDto)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "success"})
}
