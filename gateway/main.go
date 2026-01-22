package main

import (
	"gateway/middleware"
	"os"

	"github.com/bytedance/gopkg/util/logger"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {

	err := godotenv.Load("../.env", ".env")

	if err != nil {
		logger.Error("failed to load .env file", "error", err)
		os.Exit(1)
	}

	r := gin.Default()

	r.Any("/api/auth/*any", proxyToService("http://localhost:8081"))
	r.Use(middleware.JWTAuth())
	r.Any("/api/users/*any", proxyToService("http://localhost:8081"))
	r.Any("/api/ticket/*any", proxyToService("http://localhost:8082"))
	r.Any("/api/events/*any", proxyToService("http://localhost:8083"))
	r.Any("/api/notifications/*any", proxyToService("http://localhost:8084"))

	r.Run(":8000")
}
