package transport

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type TicketHandler struct {
}

func NewTicketHandler() *TicketHandler {
	return &TicketHandler{}
}

func (h *TicketHandler) RegisterRoutes(r *gin.Engine) {
	r.GET("ping", h.Ping)
}

func (h *TicketHandler) Ping(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "success"})
}
