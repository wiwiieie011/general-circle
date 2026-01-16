package transport

import (
	"errors"
	"event-service/internal/dto"
	e "event-service/internal/errors"
	"event-service/internal/services"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type CategoryHandler struct {
	service services.CategoryService
}

func NewCategoryHandler(service services.CategoryService) *CategoryHandler {
	return &CategoryHandler{service: service}
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
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "некорректный JSON"})
		return
	}

	category, err := h.service.CreateCategory(req)
	if err != nil {
		if errors.Is(err, e.ErrCategoryNameExists) {
			ctx.JSON(http.StatusConflict, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, category)
}

func (h *CategoryHandler) GetByID(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "некорректный ID"})
		return
	}

	category, err := h.service.GetCategory(uint(id))
	if err != nil {
		if errors.Is(err, e.ErrCategoryNotFound) {
			ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, category)
}

func (h *CategoryHandler) Delete(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "некорректный ID"})
		return
	}

	if err := h.service.DeleteCategory(uint(id)); err != nil {
		if errors.Is(err, e.ErrCategoryNotFound) {
			ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
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
