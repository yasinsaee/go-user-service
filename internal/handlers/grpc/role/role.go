package rolegrpc

import (
	"context"

	"github.com/yasinsaee/go-user-service/api/rolepb"
	"github.com/yasinsaee/go-user-service/internal/domain/permission"
	"github.com/yasinsaee/go-user-service/internal/domain/role"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Handler struct {
	rolepb.UnimplementedRoleServiceServer
	service    role.RoleService
	perService permission.PermissionService
}

func New(service role.RoleService, perService permission.PermissionService) *Handler {
	return &Handler{service: service, perService: perService}
}

// -- start helper
func toPermissionPB(p *permission.Permission) *rolepb.Permission {
	return &rolepb.Permission{
		Id:          p.ID.Hex(),
		Name:        p.Name,
		Description: p.Description,
	}
}

func (h *Handler) getPermissionsFromIDs(ids []primitive.ObjectID) ([]*rolepb.Permission, error) {
	var permProtos []*rolepb.Permission
	for _, id := range ids {
		p, err := h.perService.GetByID(id)
		if err != nil {
			return nil, err
		}
		permProtos = append(permProtos, toPermissionPB(p))
	}
	return permProtos, nil
}

func (h *Handler) toRolePB(r *role.Role) (*rolepb.Role, error) {
	perms, err := h.getPermissionsFromIDs(r.Permissions)
	if err != nil {
		return nil, err
	}
	return &rolepb.Role{
		Id:          r.ID.Hex(),
		Name:        r.Name,
		Description: r.Description,
		Permissions: perms,
	}, nil
}

//-- end helper

func (h *Handler) CreateRole(ctx context.Context, req *rolepb.CreateRoleRequest) (*rolepb.CreateRoleResponse, error) {
	r := &role.Role{
		Name:        req.GetName(),
		Description: req.GetDescription(),
	}

	for _, v := range req.GetPermissions() {
		id, err := primitive.ObjectIDFromHex(v)
		if err != nil {
			return nil, status.Errorf(codes.InvalidArgument, "invalid permission id: %v", err)
		}
		r.Permissions = append(r.Permissions, id)
	}

	err := h.service.Create(r)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create role: %v", err)
	}

	rolePB, err := h.toRolePB(r)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to map role: %v", err)
	}

	return &rolepb.CreateRoleResponse{
		Role: rolePB,
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
		rol.Permissions = make([]primitive.ObjectID, 0, len(req.Permissions))
		for _, v := range req.GetPermissions() {
			id, err := primitive.ObjectIDFromHex(v)
			if err != nil {
				return nil, status.Errorf(codes.InvalidArgument, "invalid permission id: %v", err)
			}
			rol.Permissions = append(rol.Permissions, id)
		}
	}

	if err := h.service.Update(rol); err != nil {
		return nil, status.Errorf(codes.Internal, "failed to update role: %v", err)
	}

	rolePB, err := h.toRolePB(rol)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to map role: %v", err)
	}

	return &rolepb.UpdateRoleResponse{Role: rolePB}, nil
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

	rolePB, err := h.toRolePB(r)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to map role: %v", err)
	}

	return &rolepb.GetRoleResponse{Role: rolePB}, nil
}

func (h *Handler) ListRoles(ctx context.Context, req *rolepb.ListRoleRequest) (*rolepb.ListRoleResponse, error) {
	roles, err := h.service.ListAll()
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to list roles: %v", err)
	}

	var pbRoles []*rolepb.Role
	for _, r := range roles {
		rolePB, err := h.toRolePB(&r)
		if err != nil {
			return nil, status.Errorf(codes.Internal, "failed to map role: %v", err)
		}
		pbRoles = append(pbRoles, rolePB)
	}

	return &rolepb.ListRoleResponse{Roles: pbRoles}, nil
}

func (h *Handler) DeleteRole(ctx context.Context, req *rolepb.DeleteRoleRequest) (*rolepb.DeleteRoleResponse, error) {
	err := h.service.Delete(req.GetId())
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "role not found: %v", err)
	}

	return &rolepb.DeleteRoleResponse{Message: "ok"}, nil
}
