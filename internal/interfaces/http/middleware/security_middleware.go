package middleware

import (
	"github.com/HasanNugroho/gin-clean/config"
	"github.com/HasanNugroho/gin-clean/pkg/errors"
	"github.com/gin-gonic/gin"
)

func SecurityMiddleware(config *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.Host != config.Security.ExpectedHost {
			c.Error(errors.ErrForbidden.WithMessage("Invalid host header"))
			c.Abort()
			c.Abort()
			return
		}
		c.Header("X-Frame-Options", config.Security.XFrameOptions)
		c.Header("Content-Security-Policy", config.Security.ContentSecurity)
		c.Header("X-XSS-Protection", config.Security.XXSSProtection)
		c.Header("Strict-Transport-Security", config.Security.StrictTransport)
		c.Header("Referrer-Policy", config.Security.ReferrerPolicy)
		c.Header("X-Content-Type-Options", config.Security.XContentTypeOpts)
		c.Header("Permissions-Policy", config.Security.PermissionsPolicy)
		c.Next()
	}
}
