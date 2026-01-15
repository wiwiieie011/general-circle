package main

import (
	"log/slog"
	"os"
	"ticket-service/internal/config"
	"ticket-service/internal/transport"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	logger := slog.New(slog.Default().Handler())

	err := godotenv.Load("../.env", ".env")

	if err != nil {
		logger.Error("failed to load .env file", "error", err)
		os.Exit(1)
	}

	db := config.Connect(logger)

	port := os.Getenv("SERVICE_PORT")
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
