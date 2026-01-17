package main

import (
	"context"
	"event-service/internal/config"
	"event-service/internal/kafka"
	"event-service/internal/models"
	"event-service/internal/repository"
	"event-service/internal/services"
	"event-service/internal/transport"
	"log/slog"
	"os"

	"github.com/robfig/cron/v3"

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

	brokers := config.KafkaBrokers()
	kafkaProducer := kafka.NewProducer(brokers, logger)
	defer kafkaProducer.Close()

	eventRepo := repository.NewEventRepository(db)
	scheduleRepo := repository.NewEventScheduleRepository(db)
	categoryRepo := repository.NewCategoryRepository(db)

	eventService := services.NewEventService(eventRepo, categoryRepo, kafkaProducer, logger)
	scheduleService := services.NewEventScheduleService(scheduleRepo, eventRepo)
	categoryService := services.NewCategoryService(categoryRepo)

	// Запустить cron для отправки напоминаний
	c := cron.New()
	_, err := c.AddFunc("0 9 * * *", func() { // Каждый день в 9:00
		ctx := context.Background()
		if err := eventService.SendEventReminders(ctx); err != nil {
			logger.Error("failed to send event reminders", "error", err)
		}
	})
	if err != nil {
		logger.Error("failed to add cron job", "error", err)
	}
	c.Start()
	defer c.Stop()

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
