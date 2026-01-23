package main

import (
	"gateway/middleware"
	"os"

	"github.com/bytedance/gopkg/util/logger"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {

	err := godotenv.Load(".env")

	if err != nil {
		logger.Warn("no .env file found, using system env", "error", err)
	}

	r := gin.Default()

	userURL := os.Getenv("USER_SERVICE_URL")
	if userURL == "" {
		userURL = "http://localhost:8081"
	}
	ticketURL := os.Getenv("TICKET_SERVICE_URL")
	if ticketURL == "" {
		ticketURL = "http://localhost:8082"
	}
	eventURL := os.Getenv("EVENT_SERVICE_URL")
	if eventURL == "" {
		eventURL = "http://localhost:8083"
	}
	notifURL := os.Getenv("NOTIFICATION_SERVICE_URL")
	if notifURL == "" {
		notifURL = "http://localhost:8084"
	}

	r.Any("/api/auth/*any", proxyToService(userURL))
	r.Use(middleware.JWTAuth())
	r.Any("/api/users/*any", proxyToService(userURL))
	r.Any("/api/ticket/*any", proxyToService(ticketURL))
	r.Any("/api/events/*any", proxyToService(eventURL))
	r.Any("/api/notifications/*any", proxyToService(notifURL))

	port := os.Getenv("PORT")
	if port == "" {
		port = "8000"
	}
	r.Run(":" + port)
}
