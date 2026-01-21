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
	eventHandler := NewEventHandler(eventService, log)
	scheduleHandler := NewEventScheduleHandler(scheduleService, log)
	categoryHandler := NewCategoryHandler(categoryService, log)

	eventHandler.RegisterRoutes(router)
	scheduleHandler.RegisterRoutes(router)
	categoryHandler.RegisterRoutes(router)
}
