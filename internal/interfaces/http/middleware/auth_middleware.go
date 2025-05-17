package middleware

import (
	"fmt"
	"strings"

	"github.com/HasanNugroho/gin-clean/internal/domain/service"
	"github.com/HasanNugroho/gin-clean/pkg/errors"
	"github.com/HasanNugroho/gin-clean/pkg/jwt"
	"github.com/HasanNugroho/gin-clean/pkg/logger"
	"github.com/gin-gonic/gin"
)

type AuthMiddleware struct {
	userService service.UserService
	logger      *logger.Logger
}

func NewAuthMiddleware(logger *logger.Logger, userService service.UserService) *AuthMiddleware {
	return &AuthMiddleware{userService: userService, logger: logger}
}

func (m *AuthMiddleware) AuthRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		fmt.Println("middleware")
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.Error(errors.Custom("UNAUTHORIZED", "missing authorization header", 401))
			c.Abort()
			return
		}

		token := strings.TrimPrefix(authHeader, "Bearer ")
		if token == authHeader {
			c.Error(errors.Custom("UNAUTHORIZED", "invalid authorization scheme", 401))
			c.Abort()
			return
		}

		claims, err := jwt.ParseToken(token)
		if err != nil {
			c.Error(errors.Custom("UNAUTHORIZED", "invalid or expired token", 401))
			c.Abort()
			return
		}

		userID, ok := claims["data"].(map[string]interface{})["user_id"].(string)
		if !ok {
			c.Error(errors.Custom("UNAUTHORIZED", "invalid token payload", 401))
			c.Abort()
			return
		}

		user, err := m.userService.GetById(c.Request.Context(), userID)
		if err != nil {
			c.Error(errors.Custom("UNAUTHORIZED", "user not found", 401))
			c.Abort()
			return
		}

		c.Set("user", user)
		c.Next()
	}
}
