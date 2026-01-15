package transport

import (
	"event-service/internal/services"
	"log/slog"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func RegisterRoutes(
	router *gin.Engine,
	log *slog.Logger,
	db *gorm.DB,
	eventService services.EventService,
) {
	eventHandler := NewEventHandler(eventService)
	eventHandler.RegisterRoutes(router)
}
