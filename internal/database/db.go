package database

import (
	"github.com/rs/zerolog/log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type PostgresStorage struct {
	db *gorm.DB
}

func NewPostgresStorage(connStr string) (*PostgresStorage, error) {
	db, err := gorm.Open(postgres.Open(connStr), &gorm.Config{})
	if err != nil {
		log.Error().Err(err).Msg("Unable to connect to database")
	}

	log.Info().Msg("Connected to PostgreSQL!")
	return &PostgresStorage{db: db}, nil
}

func (s *PostgresStorage) InitializeDatabase() (*gorm.DB, error) {

	if err := s.db.AutoMigrate(); err != nil {
		log.Error().Err(err).Msg("Failed to migrate database schema")
		return nil, err
	}

	log.Info().Msg("Database schema migrated successfully")

	// if err := s.SeedAdminUser(); err != nil {
	// 	log.Error().Err(err).Msg("Failed to seed initial admin user")
	// 	return nil, err
	// }
	return s.db, nil
}
