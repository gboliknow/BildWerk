package user

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
)

type UserHandler struct {
	userService *UserService
	logger      zerolog.Logger
}

func NewUserHandler(userService *UserService, logger zerolog.Logger) *UserHandler {
	return &UserHandler{userService: userService, logger: logger}
}

func (h *UserHandler) RegisterUserRoutes(r *gin.RouterGroup) {
	r.GET("/healthCheck", h.HandleHealthCheck)
	r.POST("/register", h.HandleRegister)
}

// Example: Health check endpoint
func (h *UserHandler) HandleHealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

// Example: User registration (connects to UserService)
func (h *UserHandler) HandleRegister(c *gin.Context) {
	// Call userService.CreateUser() here
	c.JSON(http.StatusOK, gin.H{"message": "User registered successfully"})
}
