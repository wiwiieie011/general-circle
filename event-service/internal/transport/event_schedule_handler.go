package transport

import (
	"errors"
	"event-service/internal/dto"
	"event-service/internal/services"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type EventScheduleHandler struct {
	service services.EventScheduleService
}

func NewEventScheduleHandler(service services.EventScheduleService) *EventScheduleHandler {
	return &EventScheduleHandler{service: service}
}

func (h *EventScheduleHandler) RegisterRoutes(r *gin.Engine) {
	schedules := r.Group("/schedules")
	{
		schedules.POST("", h.Create)
		schedules.GET("/:id", h.GetByID)
	}
}

func (h *EventScheduleHandler) Create(ctx *gin.Context) {
	var req dto.CreateScheduleRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "некорректный JSON"})
		return
	}

	schedule, err := h.service.CreateSchedule(req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, schedule)
}

func (h *EventScheduleHandler) GetByID(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "некорректный ID"})
		return
	}

	schedule, err := h.service.GetSchedule(uint(id))
	if err != nil {
		if errors.Is(err, dto.ErrEventScheduleNotFound) {
			ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, schedule)
}
