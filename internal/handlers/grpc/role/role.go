package rolegrpc

import (
	"context"

	"github.com/yasinsaee/go-user-service/api/rolepb"
	"github.com/yasinsaee/go-user-service/internal/domain/role"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Handler struct {
	rolepb.UnimplementedRoleServiceServer
	service role.RoleService
}

func New(service role.RoleService) *Handler {
	return &Handler{service: service}
}

func (h *Handler) CreateRole(ctx context.Context, req *rolepb.CreateRoleRequest) (*rolepb.CreateRoleResponse, error) {
	r := &role.Role{
		Name:        req.GetName(),
		Description: req.GetDescription(),
	}

	for _, v := range req.GetPermissions() {
		id, _ := primitive.ObjectIDFromHex(v)
		r.Permissions = append(r.Permissions, id)
	}

	perms := make([]string, len(r.Permissions))
	for i, id := range r.Permissions {
		perms[i] = id.Hex()
	}

	err := h.service.Create(r)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create role: %v", err)
	}

	return &rolepb.CreateRoleResponse{
		Role: &rolepb.Role{
			Id:          r.ID.Hex(),
			Name:        r.Name,
			Description: r.Description,
			Permissions: perms,
		},
	}, nil
}

func (h *Handler) UpdateRole(ctx context.Context, req *rolepb.UpdateRoleRequest) (*rolepb.UpdateRoleResponse, error) {
	id, err := primitive.ObjectIDFromHex(req.GetId())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid id format")
	}

	rol, err := h.service.GetByID(id)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "role not found: %v", err)
	}

	if name := req.GetName(); name != "" {
		rol.Name = name
	}
	if desc := req.GetDescription(); desc != "" {
		rol.Description = desc
	}
	if len(req.Permissions) > 0 {
		for _, v := range req.GetPermissions() {
			id, _ := primitive.ObjectIDFromHex(v)
			rol.Permissions = append(rol.Permissions, id)
		}
	}

	perms := make([]string, len(rol.Permissions))
	for i, id := range rol.Permissions {
		perms[i] = id.Hex()
	}

	if err := h.service.Update(rol); err != nil {
		return nil, status.Errorf(codes.Internal, "failed to update role: %v", err)
	}

	return &rolepb.UpdateRoleResponse{
		Role: &rolepb.Role{
			Id:          rol.ID.Hex(),
			Name:        rol.Name,
			Description: rol.Description,
			Permissions: perms,
		},
	}, nil
}

func (h *Handler) GetRole(ctx context.Context, req *rolepb.GetRoleRequest) (*rolepb.GetRoleResponse, error) {
	id, err := primitive.ObjectIDFromHex(req.GetId())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid id format")
	}

	r, err := h.service.GetByID(id)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "role not found: %v", err)
	}

	perms := make([]string, len(r.Permissions))
	for i, id := range r.Permissions {
		perms[i] = id.Hex()
	}

	return &rolepb.GetRoleResponse{
		Role: &rolepb.Role{
			Id:          r.ID.Hex(),
			Name:        r.Name,
			Description: r.Description,
			Permissions: perms,
		},
	}, nil
}

func (h *Handler) ListRoles(ctx context.Context, req *rolepb.ListRoleRequest) (*rolepb.ListRoleResponse, error) {
	roles, err := h.service.ListAll()
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to list roles: %v", err)
	}

	var pbRoles []*rolepb.Role
	for _, r := range roles {
		perms := make([]string, len(r.Permissions))
		for i, id := range r.Permissions {
			perms[i] = id.Hex()
		}
		pbRoles = append(pbRoles, &rolepb.Role{
			Id:          r.ID.Hex(),
			Name:        r.Name,
			Description: r.Description,
			Permissions: perms,
		})
	}

	return &rolepb.ListRoleResponse{
		Roles: pbRoles,
	}, nil
}

func (h *Handler) DeleteRole(ctx context.Context, req *rolepb.DeleteRoleRequest) (*rolepb.DeleteRoleResponse, error) {
	err := h.service.Delete(req.GetId())
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "role not found: %v", err)
	}

	return &rolepb.DeleteRoleResponse{
		Message: "ok",
	}, nil
}
