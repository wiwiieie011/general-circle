package main

import (
	"log/slog"
	"os"
	"ticket-service/internal/config"
	"ticket-service/internal/kafka"
	"ticket-service/internal/transport"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	logger := slog.New(slog.Default().Handler())

	err := godotenv.Load(".env")

	if err != nil {
		logger.Warn("no .env file found, using system env", "error", err)
	}

	db := config.DBConnect(logger)

	kafkaProducer := kafka.NewProducer(
		[]string{os.Getenv("KAFKA_BROKER")},
	)
	defer kafkaProducer.Close()

	port := os.Getenv("SERVICE_PORT")
	if port == "" {
		port = "8082"
	}

	r := gin.Default()

	transport.RegisterRoutes(r, logger, db, kafkaProducer)

	if err := r.Run(":" + port); err != nil {
		logger.Error("не удалось запустить сервер: ", slog.Any("error", err))
		os.Exit(1)
	}
}
