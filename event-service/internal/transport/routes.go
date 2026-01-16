package transport

import (
	"event-service/internal/services"
	"log/slog"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(
	router *gin.Engine,
	log *slog.Logger,
	eventService services.EventService,
	scheduleService services.EventScheduleService,
	categoryService services.CategoryService,
) {
	eventHandler := NewEventHandler(eventService)
	scheduleHandler := NewEventScheduleHandler(scheduleService)
	categoryHandler := NewCategoryHandler(categoryService)

	eventHandler.RegisterRoutes(router)
	scheduleHandler.RegisterRoutes(router)
	categoryHandler.RegisterRoutes(router)
}
