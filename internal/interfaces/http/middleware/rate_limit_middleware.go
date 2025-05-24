package middleware

import (
	"fmt"
	"net"

	"github.com/HasanNugroho/gin-clean/config"
	"github.com/HasanNugroho/gin-clean/internal/infrastructure/presistence/cache"
	"github.com/gin-gonic/gin"
	"github.com/ulule/limiter/v3"
	mgin "github.com/ulule/limiter/v3/drivers/middleware/gin"
	"github.com/ulule/limiter/v3/drivers/store/redis"
)

type RateLimit struct {
	limiter *limiter.Limiter
}

func (l *RateLimit) RateLimit() gin.HandlerFunc {
	if l.limiter == nil {
		fmt.Println("⚠️ Limiter instance is nil, skipping middleware")
		return func(c *gin.Context) {
			c.Next()
		}
	}

	fmt.Println("✅ RateLimit middleware applied")
	return mgin.NewMiddleware(l.limiter)
}

func NewRateLimiter(config *config.Config, redisClient *cache.RedisCache) (*RateLimit, error) {
	if config.Security.RateLimit == "" {
		return nil, nil
	}

	rate, err := limiter.NewRateFromFormatted(config.Security.RateLimit)
	if err != nil {
		return nil, fmt.Errorf("failed to parse rate limit format: %w", err)
	}

	ipv6Mask := net.CIDRMask(64, 128)
	options := []limiter.Option{limiter.WithIPv6Mask(ipv6Mask)}

	if redisClient == nil || redisClient.Client() == nil {
		return nil, fmt.Errorf("redis client is not initialized")
	}

	store, err := redis.NewStoreWithOptions(redisClient.Client(), limiter.StoreOptions{
		Prefix:   "limiter",
		MaxRetry: 3,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create redis store: %w", err)
	}

	return &RateLimit{
		limiter: limiter.New(store, rate, options...),
	}, nil
}
