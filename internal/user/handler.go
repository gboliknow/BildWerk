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

// HandleHealthCheck checks API health
// @Summary Health Check
// @Description Check if the API is running
// @Tags health
// @Accept  json
// @Produce  json
// @Success 200 {object} map[string]string
// @Router /healthCheck [get]
func (h *UserHandler) HandleHealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

// HandleRegister registers a new user
// @Summary Register User
// @Description Registers a new user with email and password
// @Tags user
// @Accept json
// @Produce json
// @Param request body map[string]string true "User registration details"
// @Success 201 {object} map[string]string
// @Router /register [post]
func (h *UserHandler) HandleRegister(c *gin.Context) {
	// Call userService.CreateUser() here
	c.JSON(http.StatusOK, gin.H{"message": "User registered successfully"})
}
