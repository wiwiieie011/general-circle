package transport

import (
	"log/slog"
	"os"
	api_http "ticket-service/internal/api/http"
	"ticket-service/internal/kafka"
	"ticket-service/internal/repository"
	"ticket-service/internal/services"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func RegisterRoutes(
	router *gin.Engine,
	logger *slog.Logger,
	db *gorm.DB,
	kafkaProducer *kafka.Producer,
) {
	eventClientBaseUrl := os.Getenv("EVENT_SERVICE_BASE_URL")
	if eventClientBaseUrl == "" {
		logger.Error("cannot resolve env param: EVENT_SERVICE_BASE_URL")
		os.Exit(1)
	}

	eventClient := api_http.NewEventClient(eventClientBaseUrl)

	ticketTypeRepo := repository.NewTicketTypeRepository(db)
	ticketRepo := repository.NewTicketRepository(db)

	ticketTypeService := services.NewTicketTypeService(eventClient, ticketTypeRepo)
	ticketService := services.NewTicketService(ticketRepo, ticketTypeRepo, eventClient, kafkaProducer, db, logger)

	ticketHandler := NewTicketHandler(ticketTypeService, ticketService, logger)
	ticketHandler.RegisterRoutes(router)
}
