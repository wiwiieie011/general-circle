package transport

import (
	"errors"
	"event-service/internal/dto"
	e "event-service/internal/errors"
	"event-service/internal/services"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type EventScheduleHandler struct {
	service services.EventScheduleService
	logger  *slog.Logger
}

func NewEventScheduleHandler(service services.EventScheduleService, logger *slog.Logger) *EventScheduleHandler {
	return &EventScheduleHandler{service: service, logger: logger}
}

func (h *EventScheduleHandler) RegisterRoutes(r *gin.Engine) {
	schedules := r.Group("/events/:id/schedule")
	{
		schedules.POST("", h.Create)
		schedules.GET("", h.GetByEventID)
	}
}

func (h *EventScheduleHandler) Create(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		if h.logger != nil {
			h.logger.Warn("invalid id param for create schedule", "error", err)
		}
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "некорректный ID"})
		return
	}

	var req dto.CreateScheduleRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "некорректный JSON"})
		return
	}

	schedule, err := h.service.CreateScheduleForEvent(uint(id), req)
	if err != nil {
		if errors.Is(err, e.ErrEventNotFound) {
			if h.logger != nil {
				h.logger.Warn("event not found when creating schedule", "id", id)
			}
			ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		if h.logger != nil {
			h.logger.Error("failed to create schedule", "error", err)
		}
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, schedule)
}

func (h *EventScheduleHandler) GetByEventID(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		if h.logger != nil {
			h.logger.Warn("invalid id param for get schedules", "error", err)
		}
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "некорректный ID"})
		return
	}

	schedules, err := h.service.GetScheduleByEventID(uint(id))
	if err != nil {
		if errors.Is(err, e.ErrEventNotFound) {
			if h.logger != nil {
				h.logger.Warn("event not found when getting schedules", "id", id)
			}
			ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		if h.logger != nil {
			h.logger.Error("failed to get schedules", "error", err)
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, schedules)
}
