package transport

import (
	"errors"
	"event-service/internal/dto"
	e "event-service/internal/errors"
	"event-service/internal/services"
	"log/slog"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

type EventHandler struct {
	service services.EventService
	logger  *slog.Logger
}

func NewEventHandler(service services.EventService, logger *slog.Logger) *EventHandler {
	return &EventHandler{service: service, logger: logger}
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
		if h.logger != nil {
			h.logger.Warn("invalid json for create event", "error", err)
		}
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Некорректный JSON"})
		return
	}

	event, err := h.service.CreateEvent(req)
	if err != nil {
		if h.logger != nil {
			h.logger.Error("failed to create event", "error", err)
		}
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, event)
}

func (h *EventHandler) GetByID(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		if h.logger != nil {
			h.logger.Warn("invalid id param", "error", err)
		}
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "некорректный ID"})
		return
	}

	event, err := h.service.GetEvent(uint(id))
	if err != nil {
		if errors.Is(err, e.ErrEventNotFound) {
			if h.logger != nil {
				h.logger.Warn("event not found", "id", id)
			}
			ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		if h.logger != nil {
			h.logger.Error("failed to get event", "error", err)
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, event)
}

func (h *EventHandler) Update(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		if h.logger != nil {
			h.logger.Warn("invalid id param for update", "error", err)
		}
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
		if errors.Is(err, e.ErrEventNotFound) {
			if h.logger != nil {
				h.logger.Warn("event not found for update", "id", id)
			}
			ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		if h.logger != nil {
			h.logger.Error("failed to update event", "error", err)
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, event)
}

func (h *EventHandler) Delete(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		if h.logger != nil {
			h.logger.Warn("invalid id param for delete", "error", err)
		}
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "некорректный ID"})
		return
	}

	if err := h.service.DeleteEvent(uint(id)); err != nil {
		if errors.Is(err, e.ErrEventNotFound) {
			if h.logger != nil {
				h.logger.Warn("event not found for delete", "id", id)
			}
			ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		if errors.Is(err, e.ErrEventIsNotDraft) {
			if h.logger != nil {
				h.logger.Warn("attempt to delete non-draft event", "id", id)
			}
			ctx.JSON(http.StatusConflict, gin.H{"error": err.Error()})
			return
		}
		if h.logger != nil {
			h.logger.Error("failed to delete event", "error", err)
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
		if h.logger != nil {
			h.logger.Error("failed to list events", "error", err)
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, events)
}

func (h *EventHandler) Publish(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		if h.logger != nil {
			h.logger.Warn("invalid id param for publish", "error", err)
		}
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "некорректный ID"})
		return
	}

	err = h.service.PublishEvent(uint(id))
	if err != nil {
		if errors.Is(err, e.ErrEventNotFound) {
			if h.logger != nil {
				h.logger.Warn("event not found for publish", "id", id)
			}
			ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		if errors.Is(err, e.ErrEventIsNotDraft) {
			if h.logger != nil {
				h.logger.Warn("attempt to publish non-draft event", "id", id)
			}
			ctx.JSON(http.StatusConflict, gin.H{"error": err.Error()})
			return
		}
		if h.logger != nil {
			h.logger.Error("failed to publish event", "error", err)
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "event is successfully published"})
}

func (h *EventHandler) Cancel(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		if h.logger != nil {
			h.logger.Warn("invalid id param for cancel", "error", err)
		}
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "некорректный ID"})
		return
	}

	err = h.service.CancelEvent(uint(id))
	if err != nil {
		if errors.Is(err, e.ErrEventNotFound) {
			if h.logger != nil {
				h.logger.Warn("event not found for cancel", "id", id)
			}
			ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		if errors.Is(err, e.ErrEventIsNotPublished) {
			if h.logger != nil {
				h.logger.Warn("attempt to cancel non-published event", "id", id)
			}
			ctx.JSON(http.StatusConflict, gin.H{"error": err.Error()})
			return
		}
		if h.logger != nil {
			h.logger.Error("failed to cancel event", "error", err)
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "event is cancelled"})
}

func (h *EventHandler) GetByUserID(ctx *gin.Context) {
	userID, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		if h.logger != nil {
			h.logger.Warn("invalid user id param", "error", err)
		}
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	events, err := h.service.GetEventsByUserID(uint(userID))
	if err != nil {
		if h.logger != nil {
			h.logger.Error("failed to get events by user", "error", err, "user_id", userID)
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, events)
}
