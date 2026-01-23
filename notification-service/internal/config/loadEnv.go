package config

import (
	"log/slog"

	"github.com/joho/godotenv"
)

func SetEnv(logger *slog.Logger) {
	err := godotenv.Load(".env")

	if err != nil {
		logger.Warn("no .env file found, using system env", "error", err)
	}

	logger.Info("environment variables loaded successfully")
}
