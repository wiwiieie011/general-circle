package transport

import (
	"log/slog"
	"notification-service/internal/services"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(router *gin.Engine, log *slog.Logger, notification services.NotificationService) {
	notificationHanlder:=NewNotificationHandler(notification,log)
	notificationHanlder.RegisterRoutes(router)
}
