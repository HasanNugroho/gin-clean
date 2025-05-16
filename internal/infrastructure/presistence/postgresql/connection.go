package postgresql

import (
	"fmt"
	"time"

	"github.com/HasanNugroho/gin-clean/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

const (
	maxRetries = 5
	retryDelay = 3 * time.Second
)

func NewPostgresDB(config *config.Config) (db *gorm.DB, err error) {
	dsn := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=%s",
		config.Database.User,
		config.Database.Pass,
		config.Database.Host,
		config.Database.Port,
		config.Database.Name,
		config.Database.Ssl,
	)
	for i := 1; i <= maxRetries; i++ {
		// Open database connection
		db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
			Logger: logger.Default.LogMode(logger.LogLevel(config.Server.LogLevel)),
		})

		if err == nil {
			break
		}

		if i < maxRetries {
			time.Sleep(retryDelay)
		} else {
			return nil, fmt.Errorf("❌ Failed to connect after %d attempts: %v", maxRetries, err)
		}
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get sqlDB from gorm: %w", err)
	}

	// Optional: set pooling params
	sqlDB.SetMaxOpenConns(10)
	sqlDB.SetMaxIdleConns(5)
	sqlDB.SetConnMaxLifetime(1 * time.Hour)
	return db, nil
}
