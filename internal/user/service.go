package user

import (
	"errors"
	"fmt"
	"net/http"
	"regexp"
	"sync"
	"time"

	emailing "github.com/gboliknow/bildwerk/internal/email"
	"github.com/gboliknow/bildwerk/internal/logging"
	"github.com/gboliknow/bildwerk/internal/models"
	"github.com/gboliknow/bildwerk/internal/store"
	"github.com/gboliknow/bildwerk/internal/transformer"
	"github.com/gboliknow/bildwerk/internal/utility"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
)

type OTPData struct {
	OTP       string
	ExpiresAt time.Time
}

type UserService struct {
	otpStore map[string]OTPData
	store    store.Store
	logger   zerolog.Logger
}

func NewUserService(s store.Store, logger zerolog.Logger) *UserService {
	return &UserService{store: s, logger: logger, otpStore: make(map[string]OTPData)}
}

// Business logic functions (without HTTP context)
func (s *UserService) RegisterUser(input RegisterUserDTO) (*models.UserResponse, *utility.AppError) {

	if errMsg := s.validateOTP(input.Email, input.OTP); errMsg != "" {
		log.Warn().Str("otp", input.OTP).Msg("Invalid OTP passed for email registration")
		return nil, utility.NewAppError(errMsg, http.StatusBadRequest)
	}

	user := &models.User{
		FirstName: input.FirstName,
		LastName:  input.LastName,
		Email:     input.Email,
	}

	u, err := s.store.CreateUser(user)
	if err != nil {
		if pgErr, ok := err.(*pgconn.PgError); ok && pgErr.Code == "23505" {
			return nil, utility.NewAppError("Email already exists", http.StatusConflict)
		}
		return nil, utility.NewAppError("Error creating user", http.StatusInternalServerError)
	}

	responseData := transformer.ToUserResponse(u)

	return &responseData, nil
}

func (s *UserService) SendOTP(email, subject string) error {
	otp, err := s.generateAndStoreOTP(email)
	if err != nil {
		return fmt.Errorf("error generating OTP: %w", err)
	}

	if subject == "" {
		subject = "Your OTP Code"
	}

	if _, err := emailing.SendOTPEmail(email, otp, subject); err != nil {
		return fmt.Errorf("error sending OTP email: %w", err)
	}

	return nil
}

func (s *UserService) VerifyOTP(email, otp string) error {
	if err := s.validateOTP(email, otp); err != "" {
		return fmt.Errorf(err)
	}
	return nil
}

func (s *UserService) FindUserByEmail(email string) (*models.User, *utility.AppError) {
	var user models.User
	if err := s.store.FindUserByEmail(email, &user); err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, utility.NewAppError("Email already exists", http.StatusConflict)
		}
		return nil, utility.NewAppError("Internal server error", http.StatusInternalServerError)
	}
	return &user, nil
}

func (s *UserService) handleUserLogin(input LoginRequestDTO, c *gin.Context) (*models.AuthResponse, *utility.AppError) {
	u, appErr := s.FindUserByEmail(input.Email)
	if appErr != nil {
		return nil, utility.NewAppError(appErr.Message, appErr.StatusCode)
	}

	if err := s.validateOTP(input.Email, input.OTP); err != "" {
		log.Warn().Str("otp", input.OTP).Msg("Invalid OTP passed for email registration")
		return nil, utility.NewAppError(appErr.Message, http.StatusUnauthorized)
	}

	token, err := utility.CreateAndSetAuthCookie(u.ID, c.Writer)
	if err != nil {
		return nil, utility.NewAppError("Internal server error", http.StatusInternalServerError)
	}

	responseData := transformer.ToUserResponse(u)
	var resp = models.AuthResponse{
		User:  responseData,
		Token: token,
	}

	logging.LogLoginAttempt(responseData.Email, true)
	return &resp, nil
}

var otpStore sync.Map

func (s *UserService) generateAndStoreOTP(email string) (string, error) {
	otp, err := utility.GenerateOTP()
	expiration := time.Now().Add(10 * time.Minute)
	otpStore.Store(email, OTPData{OTP: otp, ExpiresAt: expiration})
	s.logger.Info().Str("email", email).Str("otp", otp).Msg("Generated and stored OTP")

	return otp, err
}

func (s *UserService) validateOTP(email, providedOTP string) string {
	data, ok := otpStore.Load(email)
	if !ok {
		return "OTP session expired, please request a new OTP"
	}
	otpData := data.(OTPData)
	if time.Now().After(otpData.ExpiresAt) {
		return "Your OTP has expired, please request a new one"
	}
	if otpData.OTP != providedOTP {
		return "The OTP you entered is incorrect, please try again"
	}
	otpStore.Delete(email)
	return ""
}

var (
	errEmailRequired = errors.New("email is required")
	errInvalidEmail  = errors.New("invalid email format")

	errFirstNameRequired = errors.New("first name is required")
	errLastnameRequired  = errors.New("last name is required")
	errPasswordStrength  = errors.New("password must be at least 8 characters long and include at least one uppercase letter, one lowercase letter, one number, and one special character")
)

func ValidateUserPayload(user RegisterUserDTO) error {
	if user.Email == "" {
		return errEmailRequired
	}
	if !ValidateEmail(user.Email) {
		return errInvalidEmail
	}
	if user.FirstName == "" {
		return errFirstNameRequired
	}
	if user.LastName == "" {
		return errLastnameRequired
	}

	// err := ValidatePassword(user.Password)
	// if err != nil {
	// 	return err
	// }

	return nil
}

func ValidateEmail(email string) bool {
	regex := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	re := regexp.MustCompile(regex)
	return re.MatchString(email)
}
