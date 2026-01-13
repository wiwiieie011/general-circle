package transport

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type EventHandler struct {
}

func NewEventHandler() *EventHandler {
	return &EventHandler{}
}

func (h *EventHandler) RegisterRoutes(r *gin.Engine) {
	r.GET("/ping", h.Ping)
}

func (h *EventHandler) Ping(ctx *gin.Context) {
	ctx.Status(http.StatusOK)
}
