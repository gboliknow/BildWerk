package user

import (
	"net/http"

	"github.com/gboliknow/bildwerk/internal/utility"
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
	var input RegisterUserDTO
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := ValidateUserPayload(input)
	if err != nil {
		utility.RespondWithError(c, http.StatusInternalServerError, err.Error())
		return
	}

	user, err := h.userService.RegisterUser(input)
	if err != nil {
		utility.RespondWithError(c, http.StatusInternalServerError, err.Error())
		return
	}

	utility.WriteJSON(c.Writer, http.StatusCreated, "user created", user)
}
