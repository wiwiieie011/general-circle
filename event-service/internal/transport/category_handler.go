package transport

import (
	"errors"
	"event-service/internal/dto"
	"event-service/internal/services"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type CategoryHandler struct {
	service services.CategoryService
}

func NewCategoryService(service services.CategoryService) *CategoryHandler {
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
		if errors.Is(err, dto.ErrCategoryNameExists) {
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
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "некорректный JSON"})
		return
	}

	category, err := h.service.GetCategory(uint(id))
	if err != nil {
		if errors.Is(err, dto.ErrCategoryNotFound) {
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
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "некорректный JSON"})
		return
	}

	if err := h.service.DeleteCategory(uint(id)); err != nil {
		if errors.Is(err, dto.ErrCategoryNotFound) {
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
