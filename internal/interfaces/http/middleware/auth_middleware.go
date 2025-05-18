package middleware

import (
	"strings"
	"time"

	"github.com/HasanNugroho/gin-clean/internal/domain/entity"
	"github.com/HasanNugroho/gin-clean/internal/domain/service"
	"github.com/HasanNugroho/gin-clean/internal/infrastructure/presistence/cache"
	"github.com/HasanNugroho/gin-clean/pkg/errors"
	"github.com/HasanNugroho/gin-clean/pkg/jwt"
	"github.com/HasanNugroho/gin-clean/pkg/logger"
	"github.com/gin-gonic/gin"
)

type AuthMiddleware struct {
	userService service.UserService
	logger      *logger.Logger
	jwt         *jwt.TokenGenerator
	cache       *cache.RedisCache
}

func NewAuthMiddleware(logger *logger.Logger, userService service.UserService, jwt *jwt.TokenGenerator, cache *cache.RedisCache) *AuthMiddleware {
	return &AuthMiddleware{
		userService: userService,
		logger:      logger,
		jwt:         jwt,
		cache:       cache,
	}
}

func (m *AuthMiddleware) AuthRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.Error(errors.ErrUnauthorized.WithMessage("missing authorization header"))
			c.Abort()
			return
		}

		token := strings.TrimPrefix(authHeader, "Bearer ")
		if token == authHeader {
			c.Error(errors.ErrUnauthorized.WithMessage("invalid authorization scheme"))
			c.Abort()
			return
		}

		claims, err := m.jwt.ParseToken(token)
		if err != nil {
			c.Error(errors.ErrUnauthorized.WithMessage("invalid or expired token"))
			c.Abort()
			return
		}

		id, ok := claims["payload"].(string)
		if !ok {
			c.Error(errors.ErrUnauthorized.WithMessage("invalid token payload"))
			c.Abort()
			return
		}

		user := new(entity.User)
		err = m.cache.Get(c, "user:"+id, user)
		if err != nil {
			user, err = m.userService.GetById(c.Request.Context(), id)
			if err != nil {
				c.Error(errors.ErrUnauthorized.WithMessage("user not found").WithError(err))
				c.Abort()
				return
			}
			_ = m.cache.Set(c, "user:"+id, user, 2*time.Hour)
		}

		if !user.IsActive {
			c.Error(errors.ErrForbidden.WithMessage("user not active"))
			c.Abort()
			return
		}

		c.Set("user", user)
		c.Next()
	}
}
