package main

import (
	"event-service/internal/config"
	"event-service/internal/models"
	"event-service/internal/transport"
	"log/slog"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	logger := slog.New(slog.Default().Handler())

	if err := godotenv.Load(".env"); err != nil {
		logger.Error("failed to load .env file", "error", err)
		os.Exit(1)
	}

	db := config.ConnectDB(logger)

	if err := db.AutoMigrate(
		models.Event{},
		models.EventSchedule{},
		models.Category{},
	); err != nil {
		logger.Error("failed to migrate database", "error", err)
		os.Exit(1)
	}

	port := os.Getenv("PORT")

	if port == "" {
		port = "8082"
	}

	r := gin.Default()
	transport.RegisterRoutes(r, logger, db)

	if err := r.Run(":" + port); err != nil {
		logger.Error("не удалось запустить сервер: ", slog.Any("error", err))
		os.Exit(1)
	}
}
