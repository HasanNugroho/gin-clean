package handler

import (
	"net/http"

	"github.com/HasanNugroho/gin-clean/internal/interfaces/http/dto"
	"github.com/HasanNugroho/gin-clean/internal/interfaces/http/middleware"
	"github.com/HasanNugroho/gin-clean/internal/service"
	"github.com/HasanNugroho/gin-clean/pkg/errors"
	"github.com/HasanNugroho/gin-clean/pkg/logger"
	"github.com/HasanNugroho/gin-clean/pkg/response"
	"github.com/HasanNugroho/gin-clean/pkg/validation"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type AuthHandler struct {
	service  *service.AuthService
	log      *logger.Logger
	validate *validator.Validate
}

func RegisterAuthRoutes(r *gin.RouterGroup, service *service.AuthService, log *logger.Logger, validate *validator.Validate, authMiddleware *middleware.AuthMiddleware) {
	handler := NewAuthHandler(service, log, validate)
	authGroup := r.Group("v1/auth")
	{
		authGroup.POST("/login", handler.Login)
		authGroup.POST("/refresh", authMiddleware.AuthRequired(), handler.RefreshToken)
	}
	log.Info("Auth routes registered.")
}

func NewAuthHandler(service *service.AuthService, log *logger.Logger, validate *validator.Validate) *AuthHandler {
	return &AuthHandler{
		service:  service,
		log:      log,
		validate: validate,
	}
}

// Login godoc
// @Summary      User login
// @Description  Authenticate user and return access token
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        login  body  dto.LoginRequest  true  "Login credentials"
// @Success      200  {object}  response.Response{data=dto.AuthResponse}
// @Failure      400  {object}  response.Response
// @Failure      401  {object}  response.Response
// @Failure      500  {object}  response.Response
// @Router       /v1/auth/login [post]
// @Security     ApiKeyAuth
func (h *AuthHandler) Login(ctx *gin.Context) {
	req, ok := validation.ValidateBody[dto.LoginRequest](ctx, h.validate, h.log)
	if !ok {
		return
	}

	resp, err := h.service.Login(ctx.Request.Context(), *req)
	if err != nil {
		h.log.Error("Login failed", err)
		response.SendError(ctx, errors.StatusCode(err), "Login failed", err.Error())
		return
	}

	response.SendSuccess(ctx, http.StatusOK, "Login successful", resp)
}

// RefreshToken godoc
// @Summary      Refresh access token
// @Description  Use refresh token to obtain new access token
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        body  body  dto.RenewalTokenRequest  true  "Refresh token payload"
// @Success      200  {object}  response.Response{data=dto.AuthResponse}
// @Failure      400  {object}  response.Response
// @Failure      401  {object}  response.Response
// @Failure      500  {object}  response.Response
// @Router       /v1/auth/refresh [post]
// @Security     ApiKeyAuth
func (h *AuthHandler) RefreshToken(ctx *gin.Context) {
	req, ok := validation.ValidateBody[dto.RenewalTokenRequest](ctx, h.validate, h.log)
	if !ok {
		return
	}

	resp, err := h.service.RefreshToken(ctx.Request.Context(), *req)
	if err != nil {
		h.log.Error("Refresh token failed", err)
		response.SendError(ctx, errors.StatusCode(err), "Refresh token failed", err.Error())
		return
	}

	response.SendSuccess(ctx, http.StatusOK, "Token refreshed successfully", resp)
}
