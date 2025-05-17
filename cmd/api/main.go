package main

import (
	"fmt"
	"os"

	"github.com/HasanNugroho/gin-clean/config"
	"github.com/HasanNugroho/gin-clean/docs"
	"github.com/HasanNugroho/gin-clean/internal/infrastructure/di"
	"github.com/HasanNugroho/gin-clean/internal/infrastructure/presistence/postgresql"
	"github.com/HasanNugroho/gin-clean/internal/interfaces/http/middleware"
	"github.com/HasanNugroho/gin-clean/pkg/logger"
	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title           Example Rest API
// @version         1.0
// @description     This is a sample server celler server.

// @host      localhost:7000
// @BasePath  /api

// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
func main() {
	// Load configuration
	cfg, err := config.Get()
	if err != nil {
		fmt.Printf("failed to get config: %v\n", err)
		os.Exit(1)
	}

	// Initialize logger
	appLogger := logger.NewLogger(cfg.Server.LogLevel)

	// Initialize PostgreSQL connection
	db, err := postgresql.NewPostgresDB(cfg)
	if err != nil {
		appLogger.Fatal("Failed to connect to PostgreSQL", err)
	}

	// Initialize Gin engine
	engine := gin.Default()

	engine.Use(gin.Recovery())
	engine.Use(middleware.ErrorHandler(appLogger))

	// Setup Dependency Injection container with Gin engine and other dependencies
	container, err := di.Build(engine, cfg, appLogger, db)
	if err != nil {
		appLogger.Fatal("Failed to build DI container", err)
	}
	defer container.Clean()

	// Triggerhandler registration
	_ = container.Get("user-handler")
	_ = container.Get("auth-handler")
	_ = container.Get("auth-middleware")

	// Setup Swagger documentation routes
	setupSwagger(engine, cfg)

	// Start Gin HTTP server
	if err := engine.Run(fmt.Sprintf(":%s", cfg.Server.Port)); err != nil {
		fmt.Printf("Failed to run server: %v\n", err)
	}
}

func setupSwagger(r *gin.Engine, cfg *config.Config) {
	docs.SwaggerInfo.Title = cfg.Server.Name
	docs.SwaggerInfo.Description = cfg.Server.Name
	docs.SwaggerInfo.Host = fmt.Sprintf("%s:%s", cfg.Server.Host, cfg.Server.Port)
	docs.SwaggerInfo.BasePath = "/api"

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
}
