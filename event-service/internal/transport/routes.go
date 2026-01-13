package transport

import (
	"log/slog"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func RegisterRoutes(
	router *gin.Engine,
	log *slog.Logger,
	db *gorm.DB,
) {
	eventHandler := NewEventHandler()
	eventHandler.RegisterRoutes(router)
}
