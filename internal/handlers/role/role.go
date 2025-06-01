package handler

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/yasinsaee/go-user-service/internal/context"
	"github.com/yasinsaee/go-user-service/internal/domain/role"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// RoleHandler handles HTTP requests for roles
type RoleHandler struct {
	service role.RoleService
}

// NewRoleHandler creates a new RoleHandler
func NewRoleHandler(service role.RoleService) *RoleHandler {
	return &RoleHandler{service: service}
}

// RegisterRoutes registers role-related routes
func (h *RoleHandler) RegisterRoutes(e *echo.Echo) {
	g := e.Group("/roles")
	g.POST("", h.Create)
	g.GET("/:id", h.GetByID)
	g.GET("", h.ListAll)
	g.PUT("/:id", h.Update)
	g.DELETE("/:id", h.Delete)
}

// Create godoc
// @Summary Create a new role
// @Description Create a new role with given data
// @Tags roles
// @Accept json
// @Produce json
// @Param role body role.Role true "Role data"
// @Success 201 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /roles [post]
func (h *RoleHandler) Create(c echo.Context) error {
	g := c.(*context.GlobalContext)

	r := new(role.Role)
	if err := g.Bind(r); err != nil {
		return g.CreateErrorResponse(http.StatusBadRequest, err, "invalid request body")
	}

	if err := h.service.Create(r); err != nil {
		return g.CreateErrorResponse(http.StatusInternalServerError, err, "cannot save role")
	}

	return g.CreateSuccessResponse(http.StatusCreated, "role created successfully", echo.Map{
		"role": r,
	})
}

// GetByID godoc
// @Summary Get role by ID
// @Description Get a role by its ID
// @Tags roles
// @Accept json
// @Produce json
// @Param id path string true "Role ID"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Router /roles/{id} [get]
func (h *RoleHandler) GetByID(c echo.Context) error {
	g := c.(*context.GlobalContext)

	id := g.Param("id")
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return g.CreateErrorResponse(http.StatusBadRequest, err, "invalid id format")
	}

	r, err := h.service.GetByID(objID)
	if err != nil {
		return g.CreateErrorResponse(http.StatusNotFound, err, "role not found")
	}

	return g.CreateSuccessResponse(http.StatusOK, "role retrieved successfully", echo.Map{
		"role": r,
	})
}

// ListAll godoc
// @Summary List all roles
// @Description Retrieve all roles
// @Tags roles
// @Accept json
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /roles [get]
func (h *RoleHandler) ListAll(c echo.Context) error {
	g := c.(*context.GlobalContext)

	roles, err := h.service.ListAll()
	if err != nil {
		return g.CreateErrorResponse(http.StatusInternalServerError, err, "failed to list roles")
	}

	return g.CreateSuccessResponse(http.StatusOK, "roles retrieved successfully", echo.Map{
		"roles": roles,
	})
}

// Update godoc
// @Summary Update an existing role
// @Description Update role data by ID
// @Tags roles
// @Accept json
// @Produce json
// @Param id path string true "Role ID"
// @Param role body role.Role true "Role data"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /roles/{id} [put]
func (h *RoleHandler) Update(c echo.Context) error {
	g := c.(*context.GlobalContext)

	id := g.Param("id")
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return g.CreateErrorResponse(http.StatusBadRequest, err, "invalid id format")
	}

	r := new(role.Role)
	if err := g.Bind(r); err != nil {
		return g.CreateErrorResponse(http.StatusBadRequest, err, "invalid request body")
	}

	r.ID = objID

	if err := h.service.Update(r); err != nil {
		return g.CreateErrorResponse(http.StatusInternalServerError, err, "failed to update role")
	}

	return g.CreateSuccessResponse(http.StatusOK, "role updated successfully", echo.Map{
		"role": r,
	})
}

// Delete godoc
// @Summary Delete a role
// @Description Delete role by ID
// @Tags roles
// @Accept json
// @Produce json
// @Param id path string true "Role ID"
// @Success 204 {string} string "No Content"
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /roles/{id} [delete]
func (h *RoleHandler) Delete(c echo.Context) error {
	g := c.(*context.GlobalContext)

	id := g.Param("id")
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return g.CreateErrorResponse(http.StatusBadRequest, err, "invalid id format")
	}

	if err := h.service.Delete(objID); err != nil {
		return g.CreateErrorResponse(http.StatusInternalServerError, err, "failed to delete role")
	}

	return g.CreateSuccessResponse(http.StatusNoContent, "role deleted successfully", nil)
}
