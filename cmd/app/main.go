package main

import (
	"fmt"
	"os"

	"github.com/gboliknow/bildwerk/cmd/app/api"
	"github.com/gboliknow/bildwerk/internal/config"
	"github.com/gboliknow/bildwerk/internal/database"
	"github.com/gboliknow/bildwerk/internal/store"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {
	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	connStr := getDatabaseConnectionString()
	sqlStorage, err := database.NewPostgresStorage(connStr)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to connect to database")
	}
	db, err := sqlStorage.InitializeDatabase()
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to initialize database")
	}
	store := store.NewStore(db)
	port := os.Getenv("PORT")
	apiServer := api.NewAPIServer(":"+port, store)
	log.Info().Msg("Starting API server on port" + port)
	apiServer.Serve()

}

func getDatabaseConnectionString() string {
	env := os.Getenv("ENVIRONMENT")
	if env == "debug" {
		log.Info().Msg("Starting Debug BildWerk project")
		return fmt.Sprintf(
			"postgresql://%s:%s@%s/%s?sslmode=disable",
			config.Envs.DBUser,
			config.Envs.DBPassword,
			config.Envs.DBAddress,
			config.Envs.DBName,
		)
	}
	log.Info().Msg("Starting Production BildWerk project")
	connStr := os.Getenv("DB_URL")
	if connStr == "" {
		log.Fatal().Msg("DB_URL environment variable is missing")
	}
	return connStr
}
