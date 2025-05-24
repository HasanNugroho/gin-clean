package main

import (
	"fmt"
	"os"

	"github.com/HasanNugroho/gin-clean/config"
	"github.com/HasanNugroho/gin-clean/container"
	"github.com/HasanNugroho/gin-clean/docs"
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
	// Initialize Gin engine
	engine := gin.Default()

	// Build DI container
	ctn, err := container.Build(engine)
	if err != nil {
		fmt.Printf("❌ Failed to build DI container: %v\n", err)
		os.Exit(1)
	}
	defer ctn.Clean()

	var (
		rateLimit = ctn.Get("rate-limit").(*middleware.RateLimit)
		log       = ctn.Get("logger").(*logger.Logger)
		cfg       = ctn.Get("config").(*config.Config)
	)

	// Global middlware
	engine.Use(
		gin.Recovery(),
		rateLimit.RateLimit(),
		middleware.ErrorHandler(log),
		middleware.SecurityMiddleware(cfg),
	)

	// Swagger
	setupSwagger(engine, cfg)

	// Start Gin HTTP server
	if err := engine.Run(fmt.Sprintf(":%s", cfg.Server.Port)); err != nil {
		log.Fatal("❌ Failed to run server: %v\n", err)
	}
}

// setupSwagger configures the Swagger UI endpoint.
func setupSwagger(r *gin.Engine, cfg *config.Config) {
	docs.SwaggerInfo.Title = cfg.Server.Name
	docs.SwaggerInfo.Description = cfg.Server.Name
	docs.SwaggerInfo.Host = fmt.Sprintf("%s:%s", cfg.Server.Host, cfg.Server.Port)
	docs.SwaggerInfo.BasePath = "/api"

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
}
