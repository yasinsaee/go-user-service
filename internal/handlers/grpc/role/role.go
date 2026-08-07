package rolegrpc

import (
	"context"

	cContext "github.com/yasinsaee/go-user-service/internal/context"
	"github.com/yasinsaee/go-user-service/internal/context/filter"
	"github.com/yasinsaee/go-user-service/internal/domain/permission"
	"github.com/yasinsaee/go-user-service/internal/domain/role"
	rolepb "github.com/yasinsaee/go-user-service/user-service/role"
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
		Id:          p.UniqueID,
		Name:        p.Name,
		Description: p.Description,
	}
}

func (h *Handler) getPermissionsFromIDs(ids []string) ([]*rolepb.Permission, error) {
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
		Id:          r.UniqueID,
		Name:        r.Name,
		Description: r.Description,
		Key:         r.Key,
		Permissions: perms,
		Scope:       r.Scope,
	}, nil
}

//-- end helper

func (h *Handler) CreateRole(ctx context.Context, req *rolepb.CreateRoleRequest) (*rolepb.CreateRoleResponse, error) {
	r := &role.Role{
		Name:        req.GetName(),
		Description: req.GetDescription(),
		Key:         req.GetKey(),
		Scope:       req.GetScope(),
	}

	for _, v := range req.GetPermissions() {
		r.Permissions = append(r.Permissions, v)
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
	rol, err := h.service.GetByID(req.GetId())
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "role not found: %v", err)
	}

	if name := req.GetName(); name != "" {
		rol.Name = name
	}
	if desc := req.GetDescription(); desc != "" {
		rol.Description = desc
	}
	if ky := req.GetKey(); ky != "" {
		rol.Key = ky
	}
	if sc := req.GetScope(); sc != "" {
		rol.Scope = sc
	}
	if len(req.Permissions) > 0 {
		rol.Permissions = make([]string, 0, len(req.Permissions))
		for _, v := range req.GetPermissions() {
			rol.Permissions = append(rol.Permissions, v)
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
	r, err := h.service.GetByID(req.GetId())
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

func (h *Handler) ListPaginationRoles(ctx context.Context, req *rolepb.ListPaginationRolesRequest) (*rolepb.ListPaginationRolesResponse, error) {
	var metaData = cContext.MetaData{
		Limit:       int(req.GetLimit()),
		Sort:        req.GetSort(),
		CurrentPage: int(req.GetPage()),
	}

	metaDataRes, roles, err := h.service.PaginationList(metaData, role.RoleFilter{SearchFilter: filter.SearchFilter{Search: req.GetSearch()}, Scope: req.GetScope(), IsDelete: "false"})
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

	return &rolepb.ListPaginationRolesResponse{Roles: pbRoles, MetaData: &rolepb.MetaData{
		Limit:       int32(metaDataRes.Limit),
		TotalCounts: int32(metaDataRes.TotalCounts),
		TotalPages:  int32(metaDataRes.TotalPages),
		CurrentPage: int32(metaDataRes.CurrentPage),
		NextPage:    int32(metaDataRes.NextPage),
		Sort:        metaDataRes.Sort,
	}}, nil
}
