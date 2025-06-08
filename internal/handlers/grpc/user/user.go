package usergrpc

import (
	"context"

	"github.com/yasinsaee/go-user-service/api/userpb"
	"github.com/yasinsaee/go-user-service/internal/domain/role"
	"github.com/yasinsaee/go-user-service/internal/domain/user"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type Handler struct {
	userpb.UnimplementedUserServiceServer
	service  user.UserService
	rService role.RoleService
}

func New(service user.UserService, rService role.RoleService) *Handler {
	return &Handler{service: service, rService: rService}
}

func (h *Handler) Login(ctx context.Context, req *userpb.LoginRequest) (*userpb.LoginResponse, error) {
	var (
		err error
	)

	u := new(user.User)
	u, err = h.service.Login(req.GetUsername(), req.GetPassword())

	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to login user: %v", err)
	}

	return &userpb.LoginResponse{
		User: &userpb.User{
			FirstName:    u.FirstName,
			LastName:     u.LastName,
			ProfileImage: u.ProfileImage,
			Roles:        []*userpb.Role{},
			Username:     u.Username,
			Email:        u.Email,
			PhoneNumber:  u.PhoneNumber,
			IsActive:     u.IsActive,
			CreatedAt:    timestamppb.Now(),
			UpdatedAt:    timestamppb.Now(),
			LastLogin:    timestamppb.Now(),
		},
	}, nil
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
		return nil, status.Errorf(codes.Internal, "failed to login user: %v", err)
	}

	return &userpb.UserResponse{
		User: &userpb.User{
			FirstName:    u.FirstName,
			LastName:     u.LastName,
			ProfileImage: u.ProfileImage,
			Roles:        []*userpb.Role{},
			Username:     u.Username,
			Email:        u.Email,
			PhoneNumber:  u.PhoneNumber,
			IsActive:     u.IsActive,
			CreatedAt:    timestamppb.Now(),
			UpdatedAt:    timestamppb.Now(),
			LastLogin:    timestamppb.Now(),
		},
	}, nil
}
