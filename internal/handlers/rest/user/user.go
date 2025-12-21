package handler

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/yasinsaee/go-user-service/internal/context"
	"github.com/yasinsaee/go-user-service/internal/domain/user"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// UserHandler handles HTTP requests for users
type UserHandler struct {
	service user.UserService
}

// NewUserHandler creates a new UserHandler
func NewUserHandler(service user.UserService) *UserHandler {
	return &UserHandler{service: service}
}

// RegisterRoutes registers user-related routes
func (h *UserHandler) RegisterRoutes(e *echo.Echo) {
	g := e.Group("/users")
	g.POST("", h.Create)
	g.GET("/:id", h.GetByID)
	g.GET("", h.ListAll)
	g.PUT("/:id", h.Update)
	g.DELETE("/:id", h.Delete)
}

// Create godoc
// @Summary Create a new user
// @Description Create a new user with given data
// @Tags users
// @Accept json
// @Produce json
// @Param user body user.User true "User data"
// @Success 201 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /users [post]
func (h *UserHandler) Create(c echo.Context) error {
	g := c.(*context.GlobalContext)

	u := new(user.User)
	if err := g.Bind(u); err != nil {
		return g.CreateErrorResponse(http.StatusBadRequest, err, "invalid request body")
	}

	// if err := h.service.Register(u); err != nil {
	// 	return g.CreateErrorResponse(http.StatusInternalServerError, err, "cannot save user")
	// }

	return g.CreateSuccessResponse(http.StatusCreated, "user created successfully", echo.Map{
		"user": u,
	})
}

// GetByID godoc
// @Summary Get user by ID
// @Description Get a user by its ID
// @Tags users
// @Accept json
// @Produce json
// @Param id path string true "User ID"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Router /users/{id} [get]
func (h *UserHandler) GetByID(c echo.Context) error {
	g := c.(*context.GlobalContext)

	id := g.Param("id")
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return g.CreateErrorResponse(http.StatusBadRequest, err, "invalid id format")
	}

	u, err := h.service.GetByID(objID)
	if err != nil {
		return g.CreateErrorResponse(http.StatusNotFound, err, "user not found")
	}

	return g.CreateSuccessResponse(http.StatusOK, "user retrieved successfully", echo.Map{
		"user": u,
	})
}

// ListAll godoc
// @Summary List all users
// @Description Retrieve all users
// @Tags users
// @Accept json
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /users [get]
func (h *UserHandler) ListAll(c echo.Context) error {
	g := c.(*context.GlobalContext)

	users, err := h.service.ListAll()
	if err != nil {
		return g.CreateErrorResponse(http.StatusInternalServerError, err, "failed to list users")
	}

	return g.CreateSuccessResponse(http.StatusOK, "users retrieved successfully", echo.Map{
		"users": users,
	})
}

// Update godoc
// @Summary Update an existing user
// @Description Update user data by ID
// @Tags users
// @Accept json
// @Produce json
// @Param id path string true "User ID"
// @Param user body user.User true "User data"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /users/{id} [put]
func (h *UserHandler) Update(c echo.Context) error {
	g := c.(*context.GlobalContext)

	id := g.Param("id")
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return g.CreateErrorResponse(http.StatusBadRequest, err, "invalid id format")
	}

	u := new(user.User)
	if err := g.Bind(u); err != nil {
		return g.CreateErrorResponse(http.StatusBadRequest, err, "invalid request body")
	}

	u.ID = objID

	if err := h.service.Update(u); err != nil {
		return g.CreateErrorResponse(http.StatusInternalServerError, err, "failed to update user")
	}

	return g.CreateSuccessResponse(http.StatusOK, "user updated successfully", echo.Map{
		"user": u,
	})
}

// Delete godoc
// @Summary Delete a user
// @Description Delete user by ID
// @Tags users
// @Accept json
// @Produce json
// @Param id path string true "User ID"
// @Success 204 {string} string "No Content"
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /users/{id} [delete]
func (h *UserHandler) Delete(c echo.Context) error {
	g := c.(*context.GlobalContext)

	id := g.Param("id")
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return g.CreateErrorResponse(http.StatusBadRequest, err, "invalid id format")
	}

	if err := h.service.Delete(objID); err != nil {
		return g.CreateErrorResponse(http.StatusInternalServerError, err, "failed to delete user")
	}

	return g.CreateSuccessResponse(http.StatusNoContent, "user deleted successfully", nil)
}
