package usergrpc

import (
	"context"
	"sync"

	"github.com/yasinsaee/go-user-service/internal/domain/permission"
	"github.com/yasinsaee/go-user-service/internal/domain/role"
	"github.com/yasinsaee/go-user-service/internal/domain/user"
	"github.com/yasinsaee/go-user-service/pkg/jwt"
	userpb "github.com/yasinsaee/go-user-service/user-service/user"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type Handler struct {
	userpb.UnimplementedUserServiceServer
	service    user.UserService
	rService   role.RoleService
	pService   permission.PermissionService
	roleCache  map[string]*role.Role
	permCache  map[string]*permission.Permission
	cacheMutex sync.RWMutex
}

func New(service user.UserService, rService role.RoleService, pService permission.PermissionService) *Handler {
	return &Handler{
		service:   service,
		rService:  rService,
		pService:  pService,
		roleCache: make(map[string]*role.Role),
		permCache: make(map[string]*permission.Permission),
	}
}

// -- #start helpers

func toPermissionPB(p *permission.Permission) *userpb.Permission {
	return &userpb.Permission{
		Id:          p.ID.Hex(),
		Name:        p.Name,
		Description: p.Description,
	}
}

func (h *Handler) getPermissionByID(id primitive.ObjectID) (*permission.Permission, error) {
	idStr := id.Hex()

	h.cacheMutex.RLock()
	if p, ok := h.permCache[idStr]; ok {
		h.cacheMutex.RUnlock()
		return p, nil
	}
	h.cacheMutex.RUnlock()

	p, err := h.pService.GetByID(id)
	if err != nil {
		return nil, err
	}

	h.cacheMutex.Lock()
	h.permCache[idStr] = p
	h.cacheMutex.Unlock()

	return p, nil
}

func (h *Handler) getRoleByID(id primitive.ObjectID) (*role.Role, error) {
	idStr := id.Hex()

	h.cacheMutex.RLock()
	if r, ok := h.roleCache[idStr]; ok {
		h.cacheMutex.RUnlock()
		return r, nil
	}
	h.cacheMutex.RUnlock()

	r, err := h.rService.GetByID(id)
	if err != nil {
		return nil, err
	}

	h.cacheMutex.Lock()
	h.roleCache[idStr] = r
	h.cacheMutex.Unlock()

	return r, nil
}

func (h *Handler) getPermissionsFromIDs(ids []primitive.ObjectID) ([]*userpb.Permission, error) {
	var perms []*userpb.Permission
	for _, id := range ids {
		p, err := h.getPermissionByID(id)
		if err != nil {
			return nil, err
		}
		perms = append(perms, toPermissionPB(p))
	}
	return perms, nil
}

func (h *Handler) toRolePb(r *role.Role) *userpb.Role {
	perms, err := h.getPermissionsFromIDs(r.Permissions)
	if err != nil {
		return nil
	}
	return &userpb.Role{
		Id:          r.ID.Hex(),
		Name:        r.Name,
		Description: r.Description,
		Permissions: perms,
	}
}

func (h *Handler) getRoleFromIDs(ids []primitive.ObjectID) ([]*userpb.Role, error) {
	var roles []*userpb.Role
	for _, id := range ids {
		r, err := h.getRoleByID(id)
		if err != nil {
			return nil, err
		}
		roles = append(roles, h.toRolePb(r))
	}
	return roles, nil
}

func (h *Handler) toUserPb(u *user.User) *userpb.User {
	rolePb := make([]*userpb.Role, 0)
	if len(u.Roles) > 0 {
		rolePb, _ = h.getRoleFromIDs(u.Roles)

	}

	return &userpb.User{
		FirstName:    u.FirstName,
		LastName:     u.LastName,
		ProfileImage: u.ProfileImage,
		Roles:        rolePb,
		Username:     u.Username,
		Email:        u.Email,
		PhoneNumber:  u.PhoneNumber,
		IsActive:     u.IsActive,
		CreatedAt:    timestamppb.New(u.CreatedAt),
		UpdatedAt:    timestamppb.New(u.UpdatedAt),
		LastLogin:    timestamppb.New(u.LastLogin),
	}
}

func (h *Handler) toUserJwtMeta(u *user.User) (roles []string, permissions []string) {
	rolePbs, _ := h.getRoleFromIDs(u.Roles)
	for _, r := range rolePbs {
		roles = append(roles, r.Name)
		for _, p := range r.Permissions {
			permissions = append(permissions, p.Name)
		}
	}
	return
}

//-- end helpers

func (h *Handler) Login(ctx context.Context, req *userpb.LoginRequest) (*userpb.LoginResponse, error) {
	u, err := h.service.Login(req.GetUsername(), req.GetPassword())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to login user: %v", err)
	}

	roles, permissions := h.toUserJwtMeta(u)

	tokenConfig := jwt.TokenConfig{
		ID:       u.ID.Hex(),
		Username: u.Username,
		Roles:    roles,
		Access:   permissions,
	}
	accessToken, _, err := tokenConfig.GenerateAccessToken()
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to generate token: %v", err)
	}

	return &userpb.LoginResponse{
		User:        h.toUserPb(u),
		AccessToken: accessToken,
	}, nil
}

func (h *Handler) Register(ctx context.Context, req *userpb.RegisterUser) (*userpb.UserResponse, error) {
	u := &user.User{
		FirstName:    req.GetFirstName(),
		LastName:     req.GetLastName(),
		Username:     req.GetUsername(),
		ProfileImage: req.GetProfileImage(),
		Email:        req.GetEmail(),
		PhoneNumber:  req.GetPhoneNumber(),
		Password:     req.GetPassword(),
	}
	if req.GetRoles() != nil {
		for _, r := range req.GetRoles() {
			roleID, err := primitive.ObjectIDFromHex(r)
			if err != nil {
				return nil, status.Errorf(codes.InvalidArgument, "invalid role ID: %v", err)
			}
			role, err := h.getRoleByID(roleID)
			if err != nil {
				return nil, status.Errorf(codes.Internal, "failed to find role: %v", err)
			}
			u.Roles = append(u.Roles, role.ID)
		}
	}

	if err := h.service.Register(u); err != nil {
		return nil, status.Errorf(codes.Internal, "failed to register user: %v", err)
	}

	return &userpb.UserResponse{
		User: h.toUserPb(u),
	}, nil
}
