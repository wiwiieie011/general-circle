package main

import (
	"log/slog"
	"notification-service/internal/config"
	"notification-service/internal/kafka"
	"notification-service/internal/models"
	"notification-service/internal/repository"
	"notification-service/internal/services"
	"notification-service/internal/transport"
	"os"

	"github.com/gin-gonic/gin"
)

func main() {
	log := config.InitLogger()

	config.SetEnv(log)

	db := config.Connect(log)

	if err := db.AutoMigrate(
		&models.Notification{},
		&models.NotificationPreference{},
	); err != nil {
		log.Error("failed to migrate database", "error", err)
		os.Exit(1)
	}

	log.Info("migrations completed")
	notRepo := repository.NewNotificationRepo(db, log)
	notService := services.NewNotifictaonService(notRepo, log)

	consumer := kafka.NewConsumer(kafka.KafkaBrokers(), notService, log)
	go consumer.Start()


	httpServer:= gin.Default()
	transport.RegisterRoutes(
		httpServer,
		log,
		notService,
	)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	if err := httpServer.Run(":" + port); err != nil {
		log.Error("не удалось запустить сервер", slog.Any("error", err))
	}
}
