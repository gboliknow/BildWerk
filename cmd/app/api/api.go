package api

import (
	"net/http"
	"os"

	_ "github.com/gboliknow/bildwerk/docs"

	"github.com/gboliknow/bildwerk/internal/cache"
	"github.com/gboliknow/bildwerk/internal/store"
	"github.com/gboliknow/bildwerk/internal/user"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"

	swaggerFiles "github.com/swaggo/files"     // swagger embed files
	ginSwagger "github.com/swaggo/gin-swagger" // gin-swagger middleware
)

// @title BildWerk API
// @version 1.0
// @description This is the API documentation for BildWerk.
// @BasePath /api/v1
type APIServer struct {
	addr   string
	store   store.Store
	logger zerolog.Logger
	cache  *cache.RedisCache
}

func NewAPIServer(addr string, store store.Store) *APIServer {
	logger := zerolog.New(zerolog.ConsoleWriter{Out: os.Stderr}).With().Timestamp().Logger()
	redisCache := cache.NewRedisCache()
	return &APIServer{addr: addr, store: store, logger: logger,cache: redisCache}
}

func (s *APIServer) Serve() {
	router := gin.Default()
	apiV1 := router.Group("/api/v1")

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Initialize service & handler
	userService := user.NewUserService(s.store, s.logger)
	userHandler := user.NewUserHandler(userService, s.logger, s.cache)

	// Register user routes
	userHandler.RegisterUserRoutes(apiV1)

	s.logger.Info().Str("addr", s.addr).Msg("Starting API server")
	if err := http.ListenAndServe(s.addr, router); err != nil {
		s.logger.Fatal().Err(err).Msg("Server stopped")
	}
}
