package transport

import (
	"errors"
	"event-service/internal/dto"
	"event-service/internal/services"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

type EventHandler struct {
	service services.EventService
}

func NewEventHandler(service services.EventService) *EventHandler {
	return &EventHandler{service: service}
}

func (h *EventHandler) RegisterRoutes(r *gin.Engine) {
	r.GET("/ping", h.Ping)

	events := r.Group("/events")

	{
		events.GET("", h.List)
		events.POST("", h.Create)
		events.GET("/:id", h.GetByID)
		events.PUT("/:id", h.Update)
		events.DELETE("/:id", h.Delete)
		events.POST("/:id/publish", h.Publish)
		events.POST("/:id/cancel", h.Cancel)
	}

	r.GET("users/:id/events", h.GetByUserID)
}

func (h *EventHandler) Ping(ctx *gin.Context) {
	ctx.Status(http.StatusOK)
}

func (h *EventHandler) Create(ctx *gin.Context) {
	var req dto.CreateEventRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Некорректный JSON"})
		return
	}

	event, err := h.service.CreateEvent(req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, event)
}

func (h *EventHandler) GetByID(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "некорректный ID"})
		return
	}

	event, err := h.service.GetEvent(uint(id))
	if err != nil {
		if errors.Is(err, dto.ErrEventNotFound) {
			ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, event)
}

func (h *EventHandler) Update(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "некорректный ID"})
		return
	}

	var req dto.UpdateEventRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Некорректный JSON"})
		return
	}

	event, err := h.service.UpdateEvent(req, uint(id))
	if err != nil {
		if errors.Is(err, dto.ErrEventNotFound) {
			ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, event)
}

func (h *EventHandler) Delete(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "некорректный ID"})
		return
	}

	if err := h.service.DeleteEvent(uint(id)); err != nil {
		if errors.Is(err, dto.ErrEventNotFound) {
			ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		if errors.Is(err, dto.ErrEventIsNotDraft) {
			ctx.JSON(http.StatusConflict, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.Status(http.StatusOK)
}

func (h *EventHandler) List(ctx *gin.Context) {
	var query dto.EventListQuery

	if pageStr := ctx.Query("page"); pageStr != "" {
		if page, err := strconv.Atoi(strings.TrimSpace(pageStr)); err == nil {
			query.Page = page
		}
	}

	if limitStr := ctx.Query("limit"); limitStr != "" {
		if limit, err := strconv.Atoi(strings.TrimSpace(limitStr)); err == nil {
			query.Limit = limit
		}
	}

	if err := ctx.ShouldBindQuery(&query); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "некорректный параметры"})
		return
	}

	query.Title = strings.TrimSpace(query.Title)
	query.Status = strings.TrimSpace(query.Status)
	query.SortBy = strings.TrimSpace(query.SortBy)
	query.SortOrder = strings.TrimSpace(query.SortOrder)

	events, err := h.service.ListEvents(query)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, events)
}

func (h *EventHandler) Publish(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "некорректный ID"})
		return
	}

	err = h.service.PublishEvent(uint(id))
	if err != nil {
		if errors.Is(err, dto.ErrEventNotFound) {
			ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		if errors.Is(err, dto.ErrEventIsNotDraft) {
			ctx.JSON(http.StatusConflict, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "event is successfully published"})
}

func (h *EventHandler) Cancel(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "некорректный ID"})
		return
	}

	err = h.service.CancelEvent(uint(id))
	if err != nil {
		if errors.Is(err, dto.ErrEventNotFound) {
			ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		if errors.Is(err, dto.ErrEventIsNotPublished) {
			ctx.JSON(http.StatusConflict, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "event is cancelled"})
}

func (h *EventHandler) GetByUserID(ctx *gin.Context) {
	userID, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	events, err := h.service.GetEventsByUserID(uint(userID))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, events)
}
