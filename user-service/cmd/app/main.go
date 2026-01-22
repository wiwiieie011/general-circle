package main

import (
	"os"
	"time"

	"log/slog"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"user-service/internal/config"
	"user-service/internal/models"
	"user-service/internal/repository"
	"user-service/internal/services"
	"user-service/internal/transport"
	"user-service/internal/utils"
)

func main() {
	// ---------- LOGGER ----------
	log := config.InitLogger()

	// ---------- ENV ----------
	config.LoadEnv(log)

	// ---------- DB ----------
	db := config.ConnectDB(log)

	// ---------- MIGRATIONS ----------
	if err := migrate(db); err != nil {
		log.Error("failed to migrate database", "error", err)
		os.Exit(1)
	}
	log.Info("migrations completed")

	// ---------- REPOSITORIES ----------
	userRepo := repository.NewUserRepository(db)

	// ---------- TOKEN MANAGER ----------
	tokenManager := utils.NewTokenManager(
		os.Getenv("JWT_SECRET"),
		mustDuration(os.Getenv("JWT_ACCESS_TTL")),
		mustDuration(os.Getenv("JWT_REFRESH_TTL")),
		"user-service",
	)

	// ---------- SERVICES ----------
	userService := services.NewUserService(userRepo)
	authService := services.NewAuthService(userRepo, tokenManager)

	// ---------- HTTP ----------
	httpServer := gin.Default()

	userHandler := transport.NewUserHandler(
		userService,
		authService,
		log,
	)

	userHandler.RegisterRoutes(httpServer)

	// ---------- START ----------
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Info("user-service started", "port", port)

	if err := httpServer.Run(":" + port); err != nil {
		log.Error("failed to start server", slog.Any("error", err))
	}
}

// ---------- HELPERS ----------

func migrate(db *gorm.DB) error {
	return db.AutoMigrate(
		&models.User{},
	)
}

func mustDuration(value string) time.Duration {
	d, err := time.ParseDuration(value)
	if err != nil {
		panic("invalid duration: " + value)
	}
	return d
}
