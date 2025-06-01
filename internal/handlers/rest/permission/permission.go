package handler

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/yasinsaee/go-user-service/internal/context"
	"github.com/yasinsaee/go-user-service/internal/domain/permission"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// PermissionHandler handles HTTP requests for permissions
type PermissionHandler struct {
	service permission.PermissionService
}

// NewPermissionHandler creates a new PermissionHandler
func NewPermissionHandler(service permission.PermissionService) *PermissionHandler {
	return &PermissionHandler{service: service}
}

// RegisterRoutes registers permission-related routes
func (h *PermissionHandler) RegisterRoutes(e *echo.Echo) {
	g := e.Group("/permissions")
	g.POST("", h.Create)
	g.GET("/:id", h.GetByID)
	g.GET("", h.ListAll)
	g.PUT("/:id", h.Update)
	g.DELETE("/:id", h.Delete)
}

// Create godoc
// @Summary Create a new permission
// @Description Create a new permission with given data
// @Tags permissions
// @Accept json
// @Produce json
// @Param permission body permission.Permission true "Permission data"
// @Success 201 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /permissions [post]
func (h *PermissionHandler) Create(c echo.Context) error {
	g := c.(*context.GlobalContext)

	p := new(permission.Permission)
	if err := g.Bind(p); err != nil {
		return g.CreateErrorResponse(http.StatusBadRequest, err, "invalid request body")
	}

	if err := h.service.Create(p); err != nil {
		return g.CreateErrorResponse(http.StatusInternalServerError, err, "cannot save permission")
	}

	return g.CreateSuccessResponse(http.StatusCreated, "permission created successfully", echo.Map{
		"permission": p,
	})
}

// GetByID godoc
// @Summary Get permission by ID
// @Description Get a permission by its ID
// @Tags permissions
// @Accept json
// @Produce json
// @Param id path string true "Permission ID"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Router /permissions/{id} [get]
func (h *PermissionHandler) GetByID(c echo.Context) error {
	g := c.(*context.GlobalContext)

	id := g.Param("id")
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return g.CreateErrorResponse(http.StatusBadRequest, err, "invalid id format")
	}

	p, err := h.service.GetByID(objID)
	if err != nil {
		return g.CreateErrorResponse(http.StatusNotFound, err, "permission not found")
	}

	return g.CreateSuccessResponse(http.StatusOK, "permission retrieved successfully", echo.Map{
		"permission": p,
	})
}

// ListAll godoc
// @Summary List all permissions
// @Description Retrieve all permissions
// @Tags permissions
// @Accept json
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /permissions [get]
func (h *PermissionHandler) ListAll(c echo.Context) error {
	g := c.(*context.GlobalContext)

	permissions, err := h.service.ListAll()
	if err != nil {
		return g.CreateErrorResponse(http.StatusInternalServerError, err, "failed to list permissions")
	}

	return g.CreateSuccessResponse(http.StatusOK, "permissions retrieved successfully", echo.Map{
		"permissions": permissions,
	})
}

// Update godoc
// @Summary Update an existing permission
// @Description Update permission data by ID
// @Tags permissions
// @Accept json
// @Produce json
// @Param id path string true "Permission ID"
// @Param permission body permission.Permission true "Permission data"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /permissions/{id} [put]
func (h *PermissionHandler) Update(c echo.Context) error {
	g := c.(*context.GlobalContext)

	id := g.Param("id")
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return g.CreateErrorResponse(http.StatusBadRequest, err, "invalid id format")
	}

	p := new(permission.Permission)
	if err := g.Bind(p); err != nil {
		return g.CreateErrorResponse(http.StatusBadRequest, err, "invalid request body")
	}

	p.ID = objID

	if err := h.service.Update(p); err != nil {
		return g.CreateErrorResponse(http.StatusInternalServerError, err, "failed to update permission")
	}

	return g.CreateSuccessResponse(http.StatusOK, "permission updated successfully", echo.Map{
		"permission": p,
	})
}

// Delete godoc
// @Summary Delete a permission
// @Description Delete permission by ID
// @Tags permissions
// @Accept json
// @Produce json
// @Param id path string true "Permission ID"
// @Success 204 {string} string "No Content"
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /permissions/{id} [delete]
func (h *PermissionHandler) Delete(c echo.Context) error {
	g := c.(*context.GlobalContext)

	id := g.Param("id")
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return g.CreateErrorResponse(http.StatusBadRequest, err, "invalid id format")
	}

	if err := h.service.Delete(objID); err != nil {
		return g.CreateErrorResponse(http.StatusInternalServerError, err, "failed to delete permission")
	}

	return g.CreateSuccessResponse(http.StatusNoContent, "permission deleted successfully", nil)
}
