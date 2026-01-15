package transport

import (
	"log/slog"
	"os"
	api_http "ticket-service/internal/api/http"
	"ticket-service/internal/repository"
	"ticket-service/internal/services"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func RegisterRoutes(
	router *gin.Engine,
	logger *slog.Logger,
	db *gorm.DB,
) {
	eventClientBaseUrl := os.Getenv("EVENT_SERVICE_BASE_URL")
	if eventClientBaseUrl == "" {
		logger.Error("cannot resolve env param: EVENT_SERVICE_BASE_URL")
		os.Exit(1)
	}

	eventClient := api_http.NewEventClient(eventClientBaseUrl)
	ticketTypeRepo := repository.NewTicketTypeRepository(db)
	ticketTypeService := services.NewTicketTypeService(eventClient, ticketTypeRepo)
	ticketHandler := NewTicketHandler(ticketTypeService, logger)
	ticketHandler.RegisterRoutes(router)
}
