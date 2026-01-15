package transport

import (
	"log/slog"
	"net/http"
	"strconv"
	"ticket-service/internal/dto"
	"ticket-service/internal/services"

	"github.com/gin-gonic/gin"
)

type TicketHandler struct {
	ticketTypeService *services.TicketTypeService
	logger            *slog.Logger
}

func NewTicketHandler(
	ticketTypeService *services.TicketTypeService,
	logger *slog.Logger,
) *TicketHandler {
	return &TicketHandler{
		ticketTypeService: ticketTypeService,
		logger:            logger,
	}
}

func (h *TicketHandler) RegisterRoutes(r *gin.Engine) {
	r.GET("ping", h.Ping)
	r.POST("events/:id/ticket-types", h.CreateTicketType)
}

func (h *TicketHandler) Ping(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "success"})
}

func (h *TicketHandler) CreateTicketType(c *gin.Context) {
	ctx := c.Request.Context()
	idStr := c.Param("id")
	eventId, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	var ttDto dto.CreateTicketTypeRequest
	if err := c.ShouldBindJSON(&ttDto); err != nil {
		h.logger.Error(err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}
	ticketType, err := h.ticketTypeService.Create(ctx, uint(eventId), ttDto)
	if err != nil {
		h.logger.Error(err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": ticketType})
}
