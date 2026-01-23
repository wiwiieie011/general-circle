package transport

import (
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"strconv"
	e "user-service/internal/errors"
	"user-service/internal/services"
	"user-service/internal/transport/dto"
	"user-service/middleware"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	userService services.UserService
	authService *services.AuthService
	logger      *slog.Logger
}

func NewUserHandler(
	userService services.UserService,
	authService *services.AuthService,
	logger *slog.Logger,

) *UserHandler {
	return &UserHandler{
		userService: userService,
		authService: authService,
		logger:      logger,
	}
}

func (h *UserHandler) RegisterRoutes(r *gin.Engine) {
	r.GET("/ping", h.Ping)

	auth := r.Group("/auth")
	{
		auth.POST("/register", h.Register)
		auth.POST("/login", h.Login)
		auth.POST("/refresh", h.Refresh)
	}

	users := r.Group("/users")
	{
		users.GET("/me", h.GetMe)
		users.PUT("/me", h.UpdateMe)
		users.POST("/me/become-organizer", h.BecomeOrganizer)
		
		users.GET("/:id",middleware.RequireRole("organizer"), h.GetPublicProfile)
	}
}

func (h *UserHandler) Ping(ctx *gin.Context) {
	ctx.Status(http.StatusOK)
}

func (h *UserHandler) Register(ctx *gin.Context) {
	var req dto.RegisterRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "некорректный JSON"})
		return
	}

	user, access, refresh, err :=
		h.authService.Register(req.Email, req.Password, req.FirstName, req.LastName)
	if err != nil {
		if errors.Is(err, e.ErrEmailAlreadyExists) {
			ctx.JSON(http.StatusConflict, gin.H{"error": err.Error()})
			return
		}
		h.logger.Error("register failed", "error", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}

	ctx.JSON(http.StatusCreated, dto.AuthResponse{
		User:         dto.ToUserResponse(user),
		AccessToken:  access,
		RefreshToken: refresh,
	})
}

func (h *UserHandler) Login(ctx *gin.Context) {
	var req dto.LoginRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "некорректный JSON"})
		return
	}

	user, access, refresh, err :=
		h.authService.Login(req.Email, req.Password)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		return
	}

	ctx.JSON(http.StatusOK, dto.AuthResponse{
		User:         dto.ToUserResponse(user),
		AccessToken:  access,
		RefreshToken: refresh,
	})
}

func (h *UserHandler) Refresh(ctx *gin.Context) {
	var req dto.RefreshRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "некорректный JSON"})
		return
	}

	access, refresh, err :=
		h.authService.RefreshTokens(req.RefreshToken)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "invalid refresh token"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"access_token":  access,
		"refresh_token": refresh,
	})
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

func (h *UserHandler) GetMe(ctx *gin.Context) {
	userID, err := getUserID(ctx)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	user, err := h.userService.GetByIDs(userID)
	if err != nil {
		fmt.Println(err.Error())
		if errors.Is(err, e.ErrUserNotFound) {
			ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}

	ctx.JSON(http.StatusOK, dto.ToMeResponse(user))
}

func (h *UserHandler) UpdateMe(ctx *gin.Context) {
	userID, err := getUserID(ctx)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	var req dto.UpdateProfileRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "некорректный JSON"})
		return
	}

	user, err := h.userService.UpdateProfile(userID, req.FirstName, req.LastName)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}

	ctx.JSON(http.StatusOK, dto.ToMeResponse(user))
}

func (h *UserHandler) BecomeOrganizer(ctx *gin.Context) {
	userID, err := getUserID(ctx)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	user, err := h.userService.BecomeOrganizer(userID)
	if err != nil {
		if errors.Is(err, e.ErrAlreadyOrganizer) {
			ctx.JSON(http.StatusConflict, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}

	accessToken, err := h.authService.IssueAccessTokenForUser(user)
	if err != nil {
		h.logger.Error("token generation failed", "error", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "token generation failed"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"user":         dto.ToMeResponse(user),
		"access_token": accessToken,
	})

}

func (h *UserHandler) GetPublicProfile(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "некорректный ID"})
		return
	}

	user, err := h.userService.GetByID(uint(id))
	if err != nil {
		if errors.Is(err, e.ErrUserNotFound) {
			ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}

	ctx.JSON(http.StatusOK, dto.ToPublicUserResponse(user))
}
