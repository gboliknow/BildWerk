package user

import (
	"time"

	"github.com/gboliknow/bildwerk/internal/store"
	"github.com/rs/zerolog"
)

type OTPData struct {
	OTP       string
	ExpiresAt time.Time
}

type UserService struct {
	store  *store.Store
	logger zerolog.Logger
}

func NewUserService(s *store.Store, logger zerolog.Logger) *UserService {
	return &UserService{store: s, logger: logger}
}

// Business logic functions (without HTTP context)
func (s *UserService) CreateUser() error {
	// Business logic here
	return nil
}

func (s *UserService) FindUserByEmail(email string) error {
	// Business logic here
	return nil
}
