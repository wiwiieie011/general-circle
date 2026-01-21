package transport

import (
	"errors"
	"event-service/internal/dto"
	e "event-service/internal/errors"
	"event-service/internal/services"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type CategoryHandler struct {
	service services.CategoryService
	logger  *slog.Logger
}

func NewCategoryHandler(service services.CategoryService, logger *slog.Logger) *CategoryHandler {
	return &CategoryHandler{service: service, logger: logger}
}

func (h *CategoryHandler) RegisterRoutes(r *gin.Engine) {
	categories := r.Group("/categories")
	{
		categories.GET("", h.List)
		categories.POST("", h.Create)
		categories.GET("/:id", h.GetByID)
		categories.DELETE("/:id", h.Delete)
	}
}

func (h *CategoryHandler) Create(ctx *gin.Context) {
	var req dto.CreateCategoryRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		h.logger.Warn("invalid json for create category", "error", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "некорректный JSON"})
		return
	}

	category, err := h.service.CreateCategory(req)
	if err != nil {
		if errors.Is(err, e.ErrCategoryNameExists) {
			h.logger.Warn("category name exists", "name", req.Name)
			ctx.JSON(http.StatusConflict, gin.H{"error": err.Error()})
			return
		}
		h.logger.Error("failed to create category", "error", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, category)
}

func (h *CategoryHandler) GetByID(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		h.logger.Warn("invalid id param", "error", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "некорректный ID"})
		return
	}

	category, err := h.service.GetCategory(uint(id))
	if err != nil {
		if errors.Is(err, e.ErrCategoryNotFound) {
			h.logger.Warn("category not found", "id", id)
			ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		h.logger.Error("failed to get category", "error", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, category)
}

func (h *CategoryHandler) Delete(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		h.logger.Warn("invalid id param for delete", "error", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "некорректный ID"})
		return
	}

	if err := h.service.DeleteCategory(uint(id)); err != nil {
		if errors.Is(err, e.ErrCategoryNotFound) {
			h.logger.Warn("category not found for delete", "id", id)
			ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		h.logger.Error("failed to delete category", "error", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.Status(http.StatusOK)
}

func (h *CategoryHandler) List(ctx *gin.Context) {
	categories, err := h.service.ListCategories()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, categories)
}
