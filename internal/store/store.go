package store

import (
	"time"

	"github.com/gboliknow/bildwerk/internal/models"

	"github.com/google/uuid"

	"gorm.io/gorm"
)

type Store interface {
	CreateUser(user *models.User) (*models.User, error)
	FindUserByEmail(email string, user *models.User) error
	FindUserByID(userID string) (*models.User, error)
}

type Storage struct {
	db *gorm.DB
}

func NewStore(db *gorm.DB) *Storage {
	return &Storage{
		db: db,
	}
}
func (s *Storage) CreateUser(user *models.User) (*models.User, error) {
	user.ID = uuid.New().String()
	user.CreatedAt = time.Now()
	if err := s.db.Create(user).Error; err != nil {
		return nil, err
	}
	return user, nil
}


func (s *Storage) FindUserByEmail(email string, user *models.User) error {
	return s.db.Where("email = ?", email).First(user).Error
}

func (s *Storage) FindUserByID(userID string) (*models.User, error) {
	var user models.User
	if err := s.db.Where("id = ?", userID).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}
