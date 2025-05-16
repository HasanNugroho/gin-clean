package di

import (
	"time"

	"github.com/HasanNugroho/gin-clean/config"
	"github.com/HasanNugroho/gin-clean/internal/domain/repository"
	"github.com/HasanNugroho/gin-clean/internal/infrastructure/presistence/postgresql"
	"github.com/HasanNugroho/gin-clean/internal/interfaces/http/handler"
	"github.com/HasanNugroho/gin-clean/internal/service"
	"github.com/HasanNugroho/gin-clean/pkg/logger"
	"github.com/gin-gonic/gin"
	"github.com/sarulabs/di/v2"
	"gorm.io/gorm"
)

func Build(engine *gin.Engine, cfg *config.Config, log *logger.Logger, gormDB *gorm.DB) (di.Container, error) {
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
			Name: "db",
			Build: func(ctn di.Container) (interface{}, error) {
				return gormDB, nil
			},
		},
		{
			Name: "user-repository",
			Build: func(ctn di.Container) (interface{}, error) {
				return postgresql.NewUserRepository(ctn.Get("db").(*gorm.DB)), nil
			},
		},
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
			Name: "user-handler",
			Build: func(ctn di.Container) (interface{}, error) {
				handler.RegisterUserRoutes(
					ctn.Get("base-route").(*gin.RouterGroup),
					ctn.Get("user-service").(*service.UserService),
					ctn.Get("logger").(*logger.Logger),
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
