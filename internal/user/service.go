package user

import (
	"errors"
	"fmt"
	"net/http"
	"regexp"
	"sync"
	"time"
	"unicode"

	emailing "github.com/gboliknow/bildwerk/internal/email"
	"github.com/gboliknow/bildwerk/internal/models"
	"github.com/gboliknow/bildwerk/internal/store"
	"github.com/gboliknow/bildwerk/internal/utility"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
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
func (s *UserService) RegisterUser(input RegisterUserDTO) (*models.User, *utility.AppError) {

	if errMsg := s.validateOTP(input.Email, input.OTP); errMsg != "" {
		log.Warn().Str("otp", input.OTP).Msg("Invalid OTP passed for email registration")
		return nil, utility.NewAppError(errMsg, http.StatusBadRequest)
	}

	hashedPassword, err := utility.HashPassword(input.Password)
	if err != nil {
		return nil, utility.NewAppError("Error Creating User", http.StatusInternalServerError)
	}

	user := &models.User{
		Name:     input.Name,
		Email:    input.Email,
		Password: hashedPassword,
	}
	u, err := s.store.CreateUser(user)
	if err != nil {
		if pgErr, ok := err.(*pgconn.PgError); ok && pgErr.Code == "23505" {
			return nil, utility.NewAppError("Email already exists", http.StatusConflict)
		}
		return nil, utility.NewAppError("Error creating user", http.StatusInternalServerError)
	}

	return u, nil
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

func (s *UserService) FindUserByEmail(email string) error {
	// Business logic here
	return nil
}

func (s *UserService) handleUserLogin() {

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

	errPasswordRequired = errors.New("password is required")
	errPasswordStrength = errors.New("password must be at least 8 characters long and include at least one uppercase letter, one lowercase letter, one number, and one special character")
)

func ValidateUserPayload(user RegisterUserDTO) error {
	if user.Email == "" {
		return errEmailRequired
	}
	if !ValidateEmail(user.Email) {
		return errInvalidEmail
	}
	if user.Password == "" {
		return errPasswordRequired
	}

	err := ValidatePassword(user.Password)
	if err != nil {
		return err
	}

	return nil
}

func ValidateEmail(email string) bool {
	regex := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	re := regexp.MustCompile(regex)
	return re.MatchString(email)
}

func ValidatePassword(password string) error {
	if len(password) == 0 {
		return errPasswordRequired
	}

	if len(password) < 8 {
		return fmt.Errorf("password must be at least 8 characters long")
	}

	var hasUpper bool
	var hasLower bool
	var hasNumber bool
	var hasSpecial bool

	for _, char := range password {
		switch {
		case unicode.IsUpper(char):
			hasUpper = true
		case unicode.IsLower(char):
			hasLower = true
		case unicode.IsNumber(char):
			hasNumber = true
		case unicode.IsPunct(char) || unicode.IsSymbol(char):
			hasSpecial = true
		}
	}

	if !hasUpper || !hasLower || !hasNumber || !hasSpecial {
		return errPasswordStrength
	}

	return nil
}
