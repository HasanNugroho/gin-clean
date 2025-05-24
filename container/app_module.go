package container

import (
	"time"

	"github.com/HasanNugroho/gin-clean/config"
	"github.com/HasanNugroho/gin-clean/internal/domain/repository"
	"github.com/HasanNugroho/gin-clean/internal/infrastructure/presistence/postgresql"
	"github.com/HasanNugroho/gin-clean/internal/interfaces/http/handler"
	"github.com/HasanNugroho/gin-clean/internal/service"
	"github.com/HasanNugroho/gin-clean/pkg/jwt"
	"github.com/HasanNugroho/gin-clean/pkg/logger"
	"github.com/sarulabs/di/v2"
	"gorm.io/gorm"
)

func RegisterModul(builder *di.Builder) *di.Builder {
	definitions := []di.Def{
		// REPOSITORY
		{
			Name: "user-repository",
			Build: func(ctn di.Container) (interface{}, error) {
				db := ctn.Get("db").(*gorm.DB)
				return postgresql.NewUserRepository(db), nil
			},
		},

		// SERVICE
		{
			Name: "user-service",
			Build: func(ctn di.Container) (interface{}, error) {
				cfg := ctn.Get("config").(*config.Config)
				repository := ctn.Get("user-repository").(repository.UserRepository)

				return service.NewUserService(
					repository,
					time.Duration(cfg.Context.Timeout)*time.Second,
				), nil
			},
		},
		{
			Name: "auth-service",
			Build: func(ctn di.Container) (interface{}, error) {
				var (
					cfg        = ctn.Get("config").(*config.Config)
					repository = ctn.Get("user-repository").(repository.UserRepository)
					logger     = ctn.Get("logger").(*logger.Logger)
					jwt        = ctn.Get("jwt").(*jwt.TokenGenerator)
				)

				return service.NewAuthService(
					repository,
					logger,
					cfg,
					jwt,
					time.Duration(cfg.Context.Timeout)*time.Second,
				), nil
			},
		},

		// HANDLER
		{
			Name: "user-handler",
			Build: func(ctn di.Container) (interface{}, error) {
				handler.RegisterUserRoutes(&ctn)
				return nil, nil
			},
		},
		{
			Name: "auth-handler",
			Build: func(ctn di.Container) (interface{}, error) {
				handler.RegisterAuthRoutes(&ctn)
				return nil, nil
			},
		},
	}

	for _, def := range definitions {
		builder.Add(def)
	}

	return builder
}
