package di

import (
	"time"

	"github.com/HasanNugroho/gin-clean/config"
	"github.com/HasanNugroho/gin-clean/internal/domain/repository"
	"github.com/HasanNugroho/gin-clean/internal/infrastructure/presistence/cache"
	"github.com/HasanNugroho/gin-clean/internal/infrastructure/presistence/postgresql"
	"github.com/HasanNugroho/gin-clean/internal/interfaces/http/handler"
	"github.com/HasanNugroho/gin-clean/internal/interfaces/http/middleware"
	"github.com/HasanNugroho/gin-clean/internal/service"
	"github.com/HasanNugroho/gin-clean/pkg/jwt"
	"github.com/HasanNugroho/gin-clean/pkg/logger"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/sarulabs/di/v2"
	"gorm.io/gorm"
)

func Build(engine *gin.Engine, cfg *config.Config, log *logger.Logger, gormDB *gorm.DB, cache *cache.RedisCache) (di.Container, error) {
	builder, err := di.NewBuilder()
	if err != nil {
		return di.Container{}, err
	}

	definitions := []di.Def{
		{
			Name: "config",
			Build: func(ctn di.Container) (interface{}, error) {
				return cfg, nil
			},
		},
		{
			Name: "base-route",
			Build: func(ctn di.Container) (interface{}, error) {
				baseRoute := engine.Group("/api")
				return baseRoute, nil
			},
		},
		{
			Name: "logger",
			Build: func(ctn di.Container) (interface{}, error) {
				return log, nil
			},
		},
		{
			Name: "validate",
			Build: func(ctn di.Container) (interface{}, error) {
				return validator.New(), nil
			},
		},
		{
			Name: "db",
			Build: func(ctn di.Container) (interface{}, error) {
				return gormDB, nil
			},
		},
		{
			Name: "cache",
			Build: func(ctn di.Container) (interface{}, error) {
				return cache, nil
			},
		},
		{
			Name: "jwt",
			Build: func(ctn di.Container) (interface{}, error) {
				return jwt.SetJWTHelper(cfg, cache), nil
			},
		},

		// REPOSITORY
		{
			Name: "user-repository",
			Build: func(ctn di.Container) (interface{}, error) {
				return postgresql.NewUserRepository(ctn.Get("db").(*gorm.DB)), nil
			},
		},

		// SERVICE
		{
			Name: "user-service",
			Build: func(ctn di.Container) (interface{}, error) {
				return service.NewUserService(
					ctn.Get("user-repository").(repository.UserRepository),
					time.Duration(cfg.Context.Timeout)*time.Second,
				), nil
			},
		},
		{
			Name: "auth-service",
			Build: func(ctn di.Container) (interface{}, error) {
				return service.NewAuthService(
					ctn.Get("user-repository").(repository.UserRepository),
					ctn.Get("logger").(*logger.Logger),
					ctn.Get("config").(*config.Config),
					ctn.Get("jwt").(*jwt.TokenGenerator),
					time.Duration(cfg.Context.Timeout)*time.Second,
				), nil
			},
		},

		// MIDDLEWARE & DEPENDENCY
		{
			Name: "auth-middleware",
			Build: func(ctn di.Container) (interface{}, error) {
				return middleware.NewAuthMiddleware(
					ctn.Get("logger").(*logger.Logger),
					ctn.Get("user-service").(*service.UserService),
					ctn.Get("jwt").(*jwt.TokenGenerator),
					cache,
				), nil
			},
		},

		// HANDLER
		{
			Name: "user-handler",
			Build: func(ctn di.Container) (interface{}, error) {
				handler.RegisterUserRoutes(
					ctn.Get("base-route").(*gin.RouterGroup),
					ctn.Get("user-service").(*service.UserService),
					ctn.Get("logger").(*logger.Logger),
					ctn.Get("validate").(*validator.Validate),
					ctn.Get("auth-middleware").(*middleware.AuthMiddleware),
				)
				return nil, nil
			},
		},
		{
			Name: "auth-handler",
			Build: func(ctn di.Container) (interface{}, error) {
				handler.RegisterAuthRoutes(
					ctn.Get("base-route").(*gin.RouterGroup),
					ctn.Get("auth-service").(*service.AuthService),
					ctn.Get("logger").(*logger.Logger),
					ctn.Get("validate").(*validator.Validate),
					ctn.Get("auth-middleware").(*middleware.AuthMiddleware),
				)
				return nil, nil
			},
		},
	}

	for _, def := range definitions {
		builder.Add(def)
	}
	return builder.Build(), nil
}
