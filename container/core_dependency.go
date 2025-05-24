package container

import (
	"context"
	"fmt"

	"github.com/HasanNugroho/gin-clean/config"
	"github.com/HasanNugroho/gin-clean/internal/infrastructure/presistence/cache"
	"github.com/HasanNugroho/gin-clean/internal/infrastructure/presistence/postgresql"
	"github.com/HasanNugroho/gin-clean/internal/interfaces/http/middleware"
	"github.com/HasanNugroho/gin-clean/internal/service"
	"github.com/HasanNugroho/gin-clean/pkg/jwt"
	"github.com/HasanNugroho/gin-clean/pkg/logger"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/sarulabs/di/v2"
)

func RegisterDependency(builder *di.Builder) (*di.Builder, error) {
	definitions := []di.Def{
		// Config
		{
			Name: "config",
			Build: func(ctn di.Container) (interface{}, error) {
				// Load configuration
				cfg, err := config.Get()
				if err != nil {
					fmt.Printf("❌ failed to get config: %v\n", err)
					return nil, err
				}
				return cfg, nil
			},
		},

		// Logger
		{
			Name: "logger",
			Build: func(ctn di.Container) (interface{}, error) {
				// Initialize logger
				cfg := ctn.Get("config").(*config.Config)

				log := logger.NewLogger(cfg.Server.LogLevel)
				return log, nil
			},
		},

		// Validator
		{
			Name: "validate",
			Build: func(ctn di.Container) (interface{}, error) {
				return validator.New(), nil
			},
		},

		// Database
		{
			Name: "db",
			Build: func(ctn di.Container) (interface{}, error) {
				// Initialize PostgreSQL connection
				cfg := ctn.Get("config").(*config.Config)
				log := ctn.Get("logger").(*logger.Logger)

				db, err := postgresql.NewPostgresDB(cfg)
				if err != nil {
					log.Fatal("❌ Failed to connect to PostgreSQL", err)
					return nil, err
				}
				return db, nil
			},
		},

		// Cache
		{
			Name: "cache",
			Build: func(ctn di.Container) (interface{}, error) {
				// Initialize cache connection
				cfg := ctn.Get("config").(*config.Config)
				log := ctn.Get("logger").(*logger.Logger)

				e := cache.NewRedisCache(cfg)
				if err := e.Ping(context.Background()); err != nil {
					log.Fatal("❌ Failed to connect to Redis", err)
					return nil, err
				}
				return e, nil
			},
		},

		// Jwt Helper
		{
			Name: "jwt",
			Build: func(ctn di.Container) (interface{}, error) {
				cfg := ctn.Get("config").(*config.Config)
				cache := ctn.Get("cache").(*cache.RedisCache)
				return jwt.SetJWTHelper(cfg, cache), nil
			},
		},

		// Base Router
		{
			Name: "base-router",
			Build: func(ctn di.Container) (interface{}, error) {
				var (
					engine = ctn.Get("engine").(*gin.Engine)
				)

				baseRoute := engine.Group("/api")
				return baseRoute, nil
			},
		},

		// MIDDLEWARE & DEPENDENCY
		{
			Name: "auth-middleware",
			Build: func(ctn di.Container) (interface{}, error) {
				var (
					log     = ctn.Get("logger").(*logger.Logger)
					service = ctn.Get("user-service").(*service.UserService)
					jwt     = ctn.Get("jwt").(*jwt.TokenGenerator)
					cache   = ctn.Get("cache").(*cache.RedisCache)
				)
				return middleware.NewAuthMiddleware(
					log,
					service,
					jwt,
					cache,
				), nil
			},
		},
		// Initialize rate-limiter
		{
			Name: "rate-limit",
			Build: func(ctn di.Container) (interface{}, error) {
				//
				var (
					log   = ctn.Get("logger").(*logger.Logger)
					cfg   = ctn.Get("config").(*config.Config)
					cache = ctn.Get("cache").(*cache.RedisCache)
				)
				rateLimit, err := middleware.NewRateLimiter(cfg, cache)
				if err != nil {
					log.Fatal("❌ Failed to initialize rate limiter:", err)
					return nil, err
				}
				return rateLimit, nil
			},
		},
	}

	for _, def := range definitions {
		if err := builder.Add(def); err != nil {
			return nil, fmt.Errorf("failed to add dependency '%s': %w", def.Name, err)
		}
	}
	return builder, nil
}
