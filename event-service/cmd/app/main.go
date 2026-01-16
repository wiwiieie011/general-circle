package main

import (
	"event-service/internal/config"
	"event-service/internal/models"
	"event-service/internal/repository"
	"event-service/internal/services"
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
		&models.Event{},
		&models.EventSchedule{},
		&models.Category{},
	); err != nil {
		logger.Error("failed to migrate database", "error", err)
		os.Exit(1)
	}

	eventRepo := repository.NewEventRepository(db)
	scheduleRepo := repository.NewEventScheduleRepository(db)
	categoryRepo := repository.NewCategoryRepository(db)

	eventService := services.NewEventService(eventRepo, categoryRepo)
	scheduleService := services.NewEventScheduleService(scheduleRepo, eventRepo)
	categoryService := services.NewCategoryService(categoryRepo)

	r := gin.Default()
	transport.RegisterRoutes(
		r,
		logger,
		eventService,
		scheduleService,
		categoryService,
	)

	port := os.Getenv("PORT")

	if port == "" {
		port = "8083"
	}

	if err := r.Run(":" + port); err != nil {
		logger.Error("не удалось запустить сервер: ", slog.Any("error", err))
		os.Exit(1)
	}
}
