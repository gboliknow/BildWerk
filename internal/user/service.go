package user

import (
	"errors"
	"fmt"
	"regexp"
	"sync"
	"time"
	"unicode"

	"github.com/gboliknow/bildwerk/internal/models"
	"github.com/gboliknow/bildwerk/internal/store"
	"github.com/gboliknow/bildwerk/internal/utility"
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
func (s *UserService) RegisterUser(input RegisterUserDTO) (*models.User, error) {
	var existingUser models.User
	err := s.store.FindUserByEmail(input.Email, &existingUser)
	if err != nil {
		return nil, errors.New("user already exists")
	}

	if err := s.validateOTP(input.Email, input.OTP); err != "" {
		log.Warn().Str("otp", input.OTP).Msg("Invalid OTP passed for email registration")
		return nil, errors.New(err)
	}

	hashedPassword, err := utility.HashPassword(input.Password)
	if err != nil {
		return nil, err
	}

	user := &models.User{
		Name:     input.Name,
		Email:    input.Email,
		Password: hashedPassword,
	}

	return s.store.CreateUser(user)
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
