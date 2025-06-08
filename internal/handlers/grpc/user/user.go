package usergrpc

import (
	"context"

	"github.com/yasinsaee/go-user-service/api/userpb"
	"github.com/yasinsaee/go-user-service/internal/domain/permission"
	"github.com/yasinsaee/go-user-service/internal/domain/role"
	"github.com/yasinsaee/go-user-service/internal/domain/user"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type Handler struct {
	userpb.UnimplementedUserServiceServer
	service  user.UserService
	rService role.RoleService
	pService permission.PermissionService
}

func New(service user.UserService, rService role.RoleService, pService permission.PermissionService) *Handler {
	return &Handler{service: service, rService: rService, pService: pService}
}

// -- #start helper

func toPermissionPB(p *permission.Permission) *userpb.Permission {
	return &userpb.Permission{
		Id:          p.ID.Hex(),
		Name:        p.Name,
		Description: p.Description,
	}
}

func (h *Handler) getPermissionsFromIDs(ids []primitive.ObjectID) ([]*userpb.Permission, error) {
	var permProtos []*userpb.Permission
	for _, id := range ids {
		p, err := h.pService.GetByID(id)
		if err != nil {
			return nil, err
		}
		permProtos = append(permProtos, toPermissionPB(p))
	}
	return permProtos, nil
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
	var rolePr []*userpb.Role
	for _, id := range ids {
		r, err := h.rService.GetByID(id)
		if err != nil {
			return nil, err
		}
		rolePr = append(rolePr, h.toRolePb(r))
	}
	return rolePr, nil
}

func (h *Handler) toUserBp(u *user.User) *userpb.User {
	rolePb, _ := h.getRoleFromIDs(u.Roles)

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

//-- end helper

func (h *Handler) Login(ctx context.Context, req *userpb.LoginRequest) (*userpb.LoginResponse, error) {
	var (
		err error
	)

	u := new(user.User)
	u, err = h.service.Login(req.GetUsername(), req.GetPassword())

	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to login user: %v", err)
	}

	return &userpb.LoginResponse{User: h.toUserBp(u)}, nil
}

func (h *Handler) Register(ctx context.Context, req *userpb.RegisterUser) (*userpb.UserResponse, error) {
	var (
		err error
	)

	u := &user.User{
		FirstName:    req.GetFirstName(),
		LastName:     req.GetLastName(),
		Username:     req.GetUsername(),
		ProfileImage: req.GetProfileImage(),
		Email:        req.GetEmail(),
		PhoneNumber:  req.GetPhoneNumber(),
		Password:     req.GetPassword(),
	}

	for _, r := range req.GetRoles() {
		role, err := h.rService.GetByID(r)
		if err != nil {
			return nil, status.Errorf(codes.Internal, "failed to find role: %v", err)
		}
		u.Roles = append(u.Roles, role.ID)
	}

	err = h.service.Register(u)

	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to login error: %v", err)

	}

	return &userpb.UserResponse{
		User: h.toUserBp(u),
	}, nil
}
