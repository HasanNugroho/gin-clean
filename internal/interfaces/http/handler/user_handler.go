package handler

import (
	"net/http"

	"github.com/HasanNugroho/gin-clean/internal/domain/service"
	"github.com/HasanNugroho/gin-clean/internal/interfaces/http/dto"
	"github.com/HasanNugroho/gin-clean/internal/interfaces/http/middleware"
	"github.com/HasanNugroho/gin-clean/pkg/errors"
	"github.com/HasanNugroho/gin-clean/pkg/logger"
	"github.com/HasanNugroho/gin-clean/pkg/response"
	"github.com/HasanNugroho/gin-clean/pkg/validation"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type UserHandler struct {
	service  service.UserService
	log      *logger.Logger
	validate *validator.Validate
}

func RegisterUserRoutes(r *gin.RouterGroup, service service.UserService, log *logger.Logger, validate *validator.Validate, authMiddleware *middleware.AuthMiddleware) {
	handler := NewUserHandler(service, log, validate)
	userGroup := r.Group("v1/users")
	{
		userGroup.POST("", handler.Create)
		userGroup.GET("/:id", authMiddleware.AuthRequired(), handler.GetById)
		userGroup.PUT("/:id", authMiddleware.AuthRequired(), handler.Update)
		userGroup.DELETE("/:id", authMiddleware.AuthRequired(), handler.Delete)
	}
	log.Info("User routes registered.")
}

func NewUserHandler(service service.UserService, log *logger.Logger, validate *validator.Validate) *UserHandler {
	return &UserHandler{service: service, log: log, validate: validate}
}

func (h *UserHandler) validateUUID(c *gin.Context, paramName string) (string, bool) {
	id := c.Param(paramName)
	if err := h.validate.Var(id, "required,uuid"); err != nil {
		response.SendError(c, http.StatusBadRequest, "Invalid UUID", err.Error())
		return "", false
	}
	return id, true
}

// Create godoc
// @Summary      Create a user (admin)
// @Description  Admin endpoint to create a new user
// @Tags         Users
// @Accept       json
// @Produce      json
// @Param        body  body      dto.CreateUserRequest  true  "Create Request"
// @Success      201   {object}  response.Response{data=map[string]string}
// @Failure      400   {object}  response.Response
// @Failure      500   {object}  response.Response
// @Router       /v1/users [post]
func (h *UserHandler) Create(c *gin.Context) {
	req, ok := validation.ValidateBody[dto.CreateUserRequest](c, h.validate, h.log)
	if !ok {
		return
	}

	if err := h.service.Create(c.Request.Context(), req); err != nil {
		h.log.Error("Failed to create user", err)
		response.SendError(c, errors.StatusCode(err), "Failed to create user", err.Error())
		return
	}

	h.log.Info("User created", "email", req.Email)
	response.SendSuccess(c, http.StatusCreated, "User created successfully", map[string]string{"email": req.Email})
}

// GetById godoc
// @Summary      Get user by ID
// @Description  Retrieve a user by their ID
// @Tags         Users
// @Produce      json
// @Param        id   path      string  true  "User ID"
// @Success      200  {object}  response.Response{data=object}
// @Failure      404  {object}  response.Response
// @Failure      500  {object}  response.Response
// @Router       /v1/users/{id} [get]
// @Security     BearerAuth
func (h *UserHandler) GetById(c *gin.Context) {
	id, ok := h.validateUUID(c, "id")
	if !ok {
		return
	}

	user, err := h.service.GetById(c.Request.Context(), id)
	if err != nil {
		h.log.Error("Failed to get user by id", err, "id", id)
		response.SendError(c, errors.StatusCode(err), "User not found", err.Error())
		return
	}
	response.SendSuccess(c, http.StatusOK, "User fetched successfully", user)
}

// Update godoc
// @Summary      Update user by ID
// @Description  Update user data by user ID
// @Tags         Users
// @Accept       json
// @Produce      json
// @Param        id    path      string                 true  "User ID"
// @Param        user  body      dto.UpdateUserRequest  true  "User data to update"
// @Success      200   {object}  response.Response{data=map[string]string}
// @Failure      400   {object}  response.Response
// @Failure      404   {object}  response.Response
// @Failure      500   {object}  response.Response
// @Router       /v1/users/{id} [put]
// @Security     BearerAuth
func (h *UserHandler) Update(c *gin.Context) {
	id, ok := h.validateUUID(c, "id")
	if !ok {
		return
	}

	req, ok := validation.ValidateBody[dto.UpdateUserRequest](c, h.validate, h.log)
	if !ok {
		return
	}

	err := h.service.Update(c.Request.Context(), id, req)
	if err != nil {
		h.log.Error("Failed to update user", err, "user_id", id)
		response.SendError(c, errors.StatusCode(err), "Failed to update user", err.Error())
		return
	}

	h.log.Info("User updated successfully", "user_id", id)
	response.SendSuccess(c, http.StatusOK, "User updated successfully", map[string]string{"id": id})
}

// Delete godoc
// @Summary      Delete user
// @Description  Delete user by ID
// @Tags         Users
// @Produce      json
// @Param        id   path      string  true  "User ID"
// @Success      200  {object}  response.Response{data=map[string]string}
// @Failure      404  {object}  response.Response
// @Failure      500  {object}  response.Response
// @Router       /v1/users/{id} [delete]
// @Security     BearerAuth
func (h *UserHandler) Delete(c *gin.Context) {
	id, ok := h.validateUUID(c, "id")
	if !ok {
		return
	}

	if err := h.service.Delete(c.Request.Context(), id); err != nil {
		h.log.Error("Failed to delete user", err, "id", id)
		response.SendError(c, errors.StatusCode(err), "Failed to delete user", err.Error())
		return
	}

	h.log.Info("User deleted", "id", id)
	response.SendSuccess(c, http.StatusOK, "User deleted successfully", map[string]string{"id": id})
}
