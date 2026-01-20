package transport

import (
	"errors"
	"log/slog"
	"net/http"
	"notification-service/internal/dto"
	"notification-service/internal/services"
	"strconv"

	"github.com/gin-gonic/gin"
)

type NotificationHandler struct {
	srv services.NotificationService
	log *slog.Logger
}

func NewNotificationHandler(srv services.NotificationService, log *slog.Logger) *NotificationHandler {
	return &NotificationHandler{
		srv: srv,
		log: log,
	}
}

func (h *NotificationHandler) RegisterRoutes(r *gin.Engine) {
	notification := r.Group("/notifications")
	{
		notification.GET("", h.GetAllNotifications)
		notification.PUT("/:id/read", h.ReadNotificationByID)
		notification.PUT("/read-all", h.ReadAllNotification)
		notification.DELETE("/:id", h.DeleteNotification)
		notification.GET("/preferences", h.GetNotificationsPreferences)
		notification.PATCH("/preferences", h.UpdateNotificationPreferences)
		notification.GET("/unread-count", h.Count)
	}
}

func getUserID(ctx *gin.Context) (uint, error) {
	userIDStr := ctx.GetHeader("X-User-Id")
	if userIDStr == "" {
		return 0, errors.New("missing X-User-Id header")
	}

	id, err := strconv.ParseUint(userIDStr, 10, 64)
	if err != nil {
		return 0, errors.New("invalid X-User-Id header")
	}

	return uint(id), nil
}

func (h *NotificationHandler) GetAllNotifications(ctx *gin.Context) {
	userID, err := getUserID(ctx)
	if err != nil {
		h.log.Warn("unauthorized request", "error", err)
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	limitStr := ctx.DefaultQuery("limit", "50")
	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 {
		h.log.Warn("invalid limit parameter", "userID", userID, "limit", limitStr)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid limit"})
		return
	}

	var lastID uint
	if lastIDStr := ctx.Query("last_id"); lastIDStr != "" {
		val, err := strconv.ParseUint(lastIDStr, 10, 64)
		if err != nil {
			h.log.Warn("invalid last_id parameter", "userID", userID, "last_id", lastIDStr)
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid last_id"})
			return
		}
		lastID = uint(val)
	}

	list, err := h.srv.GetNotifications(ctx.Request.Context(), userID, limit, lastID)
	if err != nil {
		h.log.Error("failed to get notifications", "userID", userID, "error", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	h.log.Info("notifications returned", "userID", userID, "count", len(list))
	ctx.JSON(http.StatusOK, list)
}

func (h *NotificationHandler) ReadAllNotification(ctx *gin.Context) {
	userID, err := getUserID(ctx)
	if err != nil {
		h.log.Warn("unauthorized request", "error", err)
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	if err := h.srv.CheckAll(userID); err != nil {
		h.log.Warn("failed to mark all notifications as read", "userID", userID, "error", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	h.log.Info("all notifications marked as read", "userID", userID)
	ctx.JSON(http.StatusOK, gin.H{"message": "all notifications marked as read"})
}

func (h *NotificationHandler) ReadNotificationByID(ctx *gin.Context) {
	userID, err := getUserID(ctx)
	if err != nil {
		h.log.Warn("unauthorized request", "error", err)
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	id, err := strconv.ParseUint(ctx.Param("id"), 10, 64)
	if err != nil {
		h.log.Warn("invalid notification id", "userID", userID, "id", ctx.Param("id"))
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid notification id"})
		return
	}

	if err := h.srv.CheckNotificationsByID(userID, uint(id)); err != nil {
		h.log.Warn("failed to mark notification as read", "userID", userID, "notificationID", id, "error", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	h.log.Info("notification marked as read", "userID", userID, "notificationID", id)
	ctx.JSON(http.StatusOK, gin.H{"message": "notification marked as read"})
}

func (h *NotificationHandler) DeleteNotification(ctx *gin.Context) {
	userID, err := getUserID(ctx)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	id, err := strconv.ParseUint(ctx.Param("id"), 10, 64)
	if err != nil {
		h.log.Warn(
			"invalid notification id for delete",
			"userID", userID,
			"id", ctx.Param("id"),
		)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid notification id"})
		return
	}

	if err := h.srv.DeleteNotificationByID(userID, uint(id)); err != nil {
		h.log.Warn(
			"failed to delete notification",
			"userID", userID,
			"notificationID", id,
			"error", err,
		)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	h.log.Info(
		"notification deleted",
		"userID", userID,
		"notificationID", id,
	)
	ctx.JSON(http.StatusOK, gin.H{"message": "notification deleted"})
}

func (h *NotificationHandler) GetNotificationsPreferences(ctx *gin.Context) {
	userID, err := getUserID(ctx)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	settings, err := h.srv.GetNotificationPreferences(userID)
	if err != nil {
		h.log.Warn(
			"failed to get notification preferences",
			"userID", userID,
			"error", err,
		)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	h.log.Info(
		"notification preferences returned",
		"userID", userID,
	)
	ctx.JSON(http.StatusOK, settings)
}

func (h *NotificationHandler) UpdateNotificationPreferences(ctx *gin.Context) {

	userID, err := getUserID(ctx)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	var req dto.UpdateNotificationPreferencesRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		h.log.Warn(
			"invalid update notification preferences payload",
			"userID", userID,
			"error", err,
		)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	settings, err := h.srv.Update(userID, req)
	if err != nil {
		h.log.Warn(
			"failed to update notification preferences",
			"userID", userID,
			"error", err,
		)
		ctx.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		return
	}

	h.log.Info(
		"notification preferences updated",
		"userID", userID,
	)

	ctx.JSON(http.StatusOK, settings)
}

func (h *NotificationHandler) Count(ctx *gin.Context) {
	userID, err := getUserID(ctx)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	count, err := h.srv.Count(userID)
	if err != nil {
		h.log.Warn(
			"failed to count unread notifications",
			"userID", userID,
			"error", err,
		)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	h.log.Info(
		"unread notifications count returned",
		"userID", userID,
		"count", count,
	)
	ctx.JSON(http.StatusOK, gin.H{"unread_count": count})
}
