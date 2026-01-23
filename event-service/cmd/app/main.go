package main

import (
	"context"
	"event-service/internal/config"
	"event-service/internal/kafka"
	"event-service/internal/models"
	"event-service/internal/repository"
	"event-service/internal/services"
	"event-service/internal/transport"
	"log"
	"log/slog"
	"os"

	"github.com/robfig/cron/v3"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	logger := config.InitLogger()

	if err := godotenv.Load(".env"); err != nil {
		logger.Warn("no .env file found, using system env", "error", err)
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
	defer func() {
		if err := kafkaProducer.Close(); err != nil {
			logger.Error("failed to close kafka producer", "error", err)
		}
	}()

	eventRepo := repository.NewEventRepository(db, logger)
	scheduleRepo := repository.NewEventScheduleRepository(db, logger)
	categoryRepo := repository.NewCategoryRepository(db, logger)

	eventService := services.NewEventService(eventRepo, categoryRepo, kafkaProducer, logger)
	scheduleService := services.NewEventScheduleService(scheduleRepo, eventRepo, logger)
	categoryService := services.NewCategoryService(categoryRepo, logger)

	// Запустить cron для отправки напоминаний
	c := cron.New()
	_, err := c.AddFunc("0 9 * * *", func() { // Каждый день в 9:00
		ctx := context.Background()
		if err := eventService.SendEventReminders(ctx); err != nil {
			logger.Error("failed to send event reminders", "error", err)
		}
	})
	if err != nil {
		log.Fatal(err)
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
