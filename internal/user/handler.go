package user

import (
	"net/http"
	"time"

	"github.com/gboliknow/bildwerk/internal/cache"
	"github.com/gboliknow/bildwerk/internal/logging"
	"github.com/gboliknow/bildwerk/internal/middleware"
	"github.com/gboliknow/bildwerk/internal/utility"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
)

type UserHandler struct {
	userService *UserService
	logger      zerolog.Logger
	cache       *cache.RedisCache
}

func NewUserHandler(userService *UserService, logger zerolog.Logger, c *cache.RedisCache) *UserHandler {
	return &UserHandler{userService: userService, logger: logger, cache: c}
}

func (h *UserHandler) RegisterUserRoutes(r *gin.RouterGroup) {
	rateLimiter := middleware.RateLimitMiddleware(time.Second, 5, h.cache)
	r.GET("/healthCheck", h.HandleHealthCheck)
	r.POST("/register", h.HandleRegister)

	r.POST("/users/sendOtp", rateLimiter, h.handleSendOTP)
	r.POST("/users/verify", h.handleVerifyOTP)
}

// HandleHealthCheck checks API health
// @Summary Health Check
// @Description Check if the API is running
// @Tags health
// @Accept  json
// @Produce  json
// @Success 200 {object} map[string]string
// @Router /api/v1/healthCheck [get]
func (h *UserHandler) HandleHealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

// HandleRegister registers a new user
// @Summary Register User
// @Description Registers a new user with email and password
// @Tags Users
// @Accept json
// @Produce json
// @Param request body RegisterUserDTO true "User registration data"
// @Success 201 {object} map[string]string "User created successfully with authentication token"
// @Failure 400 {object} map[string]string "Invalid request payload"
// @Failure 401 {object} map[string]string "Invalid OTP"
// @Failure 409 {object} map[string]string "Email already exists"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /api/v1/register [post]
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

// handleSendOTP godoc
// @Summary Send OTP to user's email
// @Description This endpoint sends a one-time password (OTP) to the specified user's email for verification.
// @Tags Users
// @Accept json
// @Produce json
// @Param request body SendOTPRequestDTO true "Email Address"
// @Success 200 {object} map[string]string "OTP sent to email"
// @Failure 400 {object} map[string]string "Invalid request payload"
// @Failure 500 {object} map[string]string "Failed to send OTP"
// @Router   /api/v1/users/sendOtp [post]
func (h *UserHandler) handleSendOTP(c *gin.Context) {
	var payload SendOTPRequestDTO
	if err := c.ShouldBindJSON(&payload); err != nil {
		utility.RespondWithError(c, http.StatusBadRequest, "Invalid request payload")
		return
	}

	err := h.userService.SendOTP(payload.Email, payload.Subject)
	if err != nil {
		utility.RespondWithError(c, http.StatusInternalServerError, err.Error())
		return
	}

	utility.WriteJSON(c.Writer, http.StatusOK, "OTP has been sent to your email successfully!", nil)
}

// handleVerifyOTP godoc
// @Summary Verify OTP
// @Description This endpoint verifies the OTP sent to the user's email.
// @Tags Users
// @Accept json
// @Produce json
// @Param request body VerifyOTPRequestDTO true "Email and OTP"
// @Success 200 {object} map[string]string "OTP verified"
// @Failure 400 {object} map[string]string "Invalid request payload"
// @Failure 401 {object} map[string]string "OTP verification failed"
// @Router   /api/v1/users/verify [post]
func (h *UserHandler) handleVerifyOTP(c *gin.Context) {
	var payload VerifyOTPRequestDTO
	if err := c.ShouldBindJSON(&payload); err != nil {
		utility.RespondWithError(c, http.StatusBadRequest, "Invalid request payload")
		return
	}

	err := h.userService.VerifyOTP(payload.Email, payload.OTP)
	if err != nil {
		utility.RespondWithError(c, http.StatusUnauthorized, err.Error())
		logging.LogOtpVerification(payload.Email, false)
		return
	}

	logging.LogOtpVerification(payload.Email, true)
	utility.WriteJSON(c.Writer, http.StatusOK, "OTP verified", nil)
}
