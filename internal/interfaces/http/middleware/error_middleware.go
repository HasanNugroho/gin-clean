package middleware

import (
	"errors"
	"net/http"

	customError "github.com/HasanNugroho/gin-clean/pkg/errors"
	"github.com/HasanNugroho/gin-clean/pkg/logger"
	"github.com/HasanNugroho/gin-clean/pkg/response"
	"github.com/gin-gonic/gin"
)

func ErrorHandler(log *logger.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		err := c.Errors.Last()
		if err != nil {
			log.Error("Request error: %v", err.Err)

			var appErr *customError.AppError
			if errors.As(err.Err, &appErr) {
				response.SendError(c, appErr.Status, appErr.Code, appErr.Message)
				return
			}

			response.SendError(c, http.StatusInternalServerError, "INTERNAL_SERVER_ERROR", "Internal server error")
			return
		}
	}
}
