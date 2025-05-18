package main

import (
	"context"
	"fmt"
	"os"

	"github.com/HasanNugroho/gin-clean/config"
	"github.com/HasanNugroho/gin-clean/docs"
	"github.com/HasanNugroho/gin-clean/internal/infrastructure/di"
	"github.com/HasanNugroho/gin-clean/internal/infrastructure/presistence/cache"
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

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
func main() {
	// Load configuration
	cfg, err := config.Get()
	if err != nil {
		fmt.Printf("❌ failed to get config: %v\n", err)
		os.Exit(1)
	}

	// Initialize logger
	log := logger.NewLogger(cfg.Server.LogLevel)

	// Initialize PostgreSQL connection
	db, err := postgresql.NewPostgresDB(cfg)
	if err != nil {
		log.Fatal("❌ Failed to connect to PostgreSQL", err)
	}

	// Initialize cache connection
	e := cache.NewRedisCache(cfg)
	if err := e.Ping(context.Background()); err != nil {
		log.Fatal("❌ Failed to connect to Redis", err)
	}

	// Initialize rate-limit
	rateLimit, err := middleware.NewRateLimiter(cfg, e)
	if err != nil {
		log.Fatal("❌ Failed to initialize rate limiter:", err)
	}

	// Initialize Gin engine
	engine := gin.Default()
	engine.Use(
		gin.Recovery(),
		rateLimit.RateLimit(),
		middleware.ErrorHandler(log),
		middleware.SecurityMiddleware(cfg),
	)

	// Setup Dependency Injection container with Gin engine and other dependencies
	container, err := di.Build(engine, cfg, log, db, e)
	if err != nil {
		log.Fatal("❌ Failed to build DI container", err)
	}
	defer container.Clean()

	// Triggerhandler registration
	_ = container.Get("user-handler")
	_ = container.Get("auth-handler")
	_ = container.Get("auth-middleware")

	// Swagger
	setupSwagger(engine, cfg)

	// Start Gin HTTP server
	if err := engine.Run(fmt.Sprintf(":%s", cfg.Server.Port)); err != nil {
		log.Fatal("❌ Failed to run server: %v\n", err)
	}
}

func setupSwagger(r *gin.Engine, cfg *config.Config) {
	docs.SwaggerInfo.Title = cfg.Server.Name
	docs.SwaggerInfo.Description = cfg.Server.Name
	docs.SwaggerInfo.Host = fmt.Sprintf("%s:%s", cfg.Server.Host, cfg.Server.Port)
	docs.SwaggerInfo.BasePath = "/api"

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
}
