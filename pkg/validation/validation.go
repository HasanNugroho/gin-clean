package validation

import (
	"net/http"

	"github.com/HasanNugroho/gin-clean/pkg/logger"
	"github.com/HasanNugroho/gin-clean/pkg/response"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

func ValidateBody[T any](c *gin.Context, v *validator.Validate, log *logger.Logger) (*T, bool) {
	var body T

	if err := c.ShouldBindJSON(&body); err != nil {
		log.Warn("Invalid request payload", "error", err)
		response.SendError(c, http.StatusBadRequest, "Invalid request payload", err.Error())
		return nil, false
	}

	if err := v.Struct(body); err != nil {
		response.SendError(c, http.StatusBadRequest, "Validation error", err.Error())
		return nil, false
	}

	return &body, true
}
