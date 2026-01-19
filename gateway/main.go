package main

import (
	"gateway/middleware"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	r.Use(middleware.JWTAuth())
	r.Any("/api/user/*any", proxyToService("http://localhost:8081"))
	r.Any("/api/ticket/*any", proxyToService("http://localhost:8082"))
	r.Any("/api/events/*any", proxyToService("http://localhost:8083"))
	r.Any("/api/notifications/*any", proxyToService("http://localhost:8084"))

	r.Run(":8000")
}