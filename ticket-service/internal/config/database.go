package config

import (
	"fmt"
	"log/slog"
	"os"
	"ticket-service/internal/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func DBConnect(logger *slog.Logger) *gorm.DB {
	dbHost := os.Getenv("POSTGRES_HOST")
	dbUser := os.Getenv("POSTGRES_USER")
	dbPass := os.Getenv("POSTGRES_PASSWORD")
	dbName := os.Getenv("TICKETS_DB_NAME")
	dbPort := os.Getenv("POSTGRES_PORT")
	dbMode := os.Getenv("DB_SSLMODE")

	dsn := fmt.Sprintf("host=%v user=%v password=%v dbname=%v port=%v sslmode=%v",
		dbHost, dbUser, dbPass, dbName, dbPort, dbMode)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		logger.Error("failed to connect", "error", err)
		os.Exit(1)
	}

	if err := db.AutoMigrate(
		&models.TicketType{},
		&models.Ticket{},
	); err != nil {
		logger.Error("failed to migrate database", "error", err)
		os.Exit(1)
	}

	return db
}
