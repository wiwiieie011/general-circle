package transport

import (
	"log/slog"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func RegisterRoutes(
	router *gin.Engine,
	logger *slog.Logger,
	db *gorm.DB,
) {
	ticketHandler := NewTicketHandler()
	ticketHandler.RegisterRoutes(router)
}
