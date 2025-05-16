package handler

import (
	"net/http"

	"github.com/HasanNugroho/gin-clean/internal/domain/service"
	"github.com/HasanNugroho/gin-clean/internal/interfaces/http/dto"
	"github.com/HasanNugroho/gin-clean/pkg/errors"
	"github.com/HasanNugroho/gin-clean/pkg/logger"
	"github.com/HasanNugroho/gin-clean/pkg/response"
	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	service service.UserService
	log     *logger.Logger
}

func RegisterUserRoutes(r *gin.RouterGroup, service service.UserService, log *logger.Logger) {
	handler := NewUserHandler(service, log)
	userGroup := r.Group("v1/users")
	{
		userGroup.POST("", handler.Create)
		userGroup.GET("/:id", handler.GetById)
		userGroup.PUT("/:id", handler.Update)
		userGroup.DELETE("/:id", handler.Delete)
	}
	log.Info("User routes registered.")
}

func NewUserHandler(service service.UserService, log *logger.Logger) *UserHandler {
	return &UserHandler{service: service, log: log}
}

// Create godoc
// @Summary      Create a user (admin)
// @Description  Admin endpoint to create a new user
// @Tags         Users
// @Accept       json
// @Produce      json
// @Param        body  body      dto.CreateUserRequest  true  "Create Request"
// @Success      201   {object}  map[string]string
// @Failure      400   {object}  map[string]string
// @Failure      500   {object}  map[string]string
// @Router       /v1/users [post]
func (h *UserHandler) Create(c *gin.Context) {
	var req dto.CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.log.Warn("Invalid create user request", "error", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	if err := h.service.Create(c.Request.Context(), &req); err != nil {
		h.log.Error("Failed to create user", err)
		c.JSON(errors.StatusCode(err), gin.H{"error": err.Error()})
		return
	}

	h.log.Info("User created", "email", req.Email)
	c.JSON(http.StatusCreated, gin.H{"message": "user created successfully"})
}

// GetById godoc
// @Summary      Get user by ID
// @Description  Retrieve a user by their ID
// @Tags         Users
// @Produce      json
// @Param        id   path      string  true  "User ID"
// @Success      200  {object}  map[string]interface{}
// @Failure      404  {object}  map[string]string
// @Router       /v1/users/{id} [get]
func (h *UserHandler) GetById(c *gin.Context) {
	id := c.Param("id")
	user, err := h.service.GetById(c.Request.Context(), id)
	if err != nil {
		h.log.Error("Failed to get user by id", err, "id", id)
		c.JSON(errors.StatusCode(err), gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, user)
}

// Update godoc
// @Summary Update user by ID
// @Description Update user data by user ID
// @Tags Users
// @Accept json
// @Produce json
// @Param id path string true "User ID"
// @Param user body dto.UpdateUserRequest true "User data to update"
// @Success 200 {object} response.Response{data=object} "User updated successfully"
// @Failure 400 {object} response.Response "Invalid request payload"
// @Failure 404 {object} response.Response "User not found"
// @Failure 500 {object} response.Response "Internal server error"
// @Router /v1/users/{id} [put]
func (h *UserHandler) Update(c *gin.Context) {
	var req dto.UpdateUserRequest
	userID := c.Param("id")

	if err := c.ShouldBindJSON(&req); err != nil {
		h.log.Warn("Invalid update request", "error", err)
		c.JSON(http.StatusBadRequest, response.Err("invalid_request", "Invalid request payload"))
		return
	}

	err := h.service.Update(c.Request.Context(), userID, &req)
	if err != nil {
		h.log.Error("Failed to update user", err, "user_id", userID)
		statusCode := errors.StatusCode(err)
		msg := err.Error()
		if statusCode == http.StatusNotFound {
			c.JSON(statusCode, response.Err("not_found", msg))
		} else {
			c.JSON(statusCode, response.Err("internal_error", msg))
		}
		return
	}

	h.log.Info("User updated successfully", "user_id", userID)
	c.JSON(http.StatusOK, response.Ok(map[string]string{"message": "User updated successfully"}))
}

// Delete godoc
// @Summary      Delete user
// @Description  Delete user by ID
// @Tags         Users
// @Param        id   path      string  true  "User ID"
// @Success      200  {object}  map[string]string
// @Failure      404  {object}  map[string]string
// @Router       /v1/users/{id} [delete]
func (h *UserHandler) Delete(c *gin.Context) {
	id := c.Param("id")
	if err := h.service.Delete(c.Request.Context(), id); err != nil {
		h.log.Error("Failed to delete user", err, "id", id)
		c.JSON(errors.StatusCode(err), gin.H{"error": err.Error()})
		return
	}
	h.log.Info("User deleted", "id", id)
	c.JSON(http.StatusOK, gin.H{"message": "user deleted successfully"})
}
