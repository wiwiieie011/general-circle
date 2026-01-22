package config

import (
	"log/slog"

	"github.com/joho/godotenv"
)

func LoadEnv(logger *slog.Logger) {
	if err := godotenv.Load("../.env", ".env"); err != nil {
		logger.Warn("no .env file found, using system env")
		return
	}

	logger.Info("environment variables loaded")
}
