package usergrpc

import (
	"context"
	"sync"

	"github.com/yasinsaee/go-user-service/internal/app/config"
	cContext "github.com/yasinsaee/go-user-service/internal/context"
	"github.com/yasinsaee/go-user-service/internal/context/filter"
	"github.com/yasinsaee/go-user-service/pkg/util"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/yasinsaee/go-user-service/internal/domain/group"
	"github.com/yasinsaee/go-user-service/internal/domain/permission"
	"github.com/yasinsaee/go-user-service/internal/domain/role"
	"github.com/yasinsaee/go-user-service/internal/domain/user"
	userpb "github.com/yasinsaee/go-user-service/user-service/user"
)

type Handler struct {
	userpb.UnimplementedUserServiceServer
	service    user.UserService
	rService   role.RoleService
	pService   permission.PermissionService
	gService   group.GroupService
	roleCache  map[string]*role.Role
	permCache  map[string]*permission.Permission
	cacheMutex sync.RWMutex
}

func New(service user.UserService, rService role.RoleService, pService permission.PermissionService, gService group.GroupService) *Handler {
	return &Handler{
		service:   service,
		rService:  rService,
		pService:  pService,
		roleCache: make(map[string]*role.Role),
		permCache: make(map[string]*permission.Permission),
		gService:  gService,
	}
}

func (h *Handler) Update(ctx context.Context, req *userpb.UpdateUser) (*userpb.UserResponse, error) {
	var (
		err error
	)
	lType := config.GetEnv("LOGIN_TYPE", "phone")
	u := new(user.User)

	if req.GetId() == "" {
		switch lType {
		case "phone":
			u, err = h.service.GetByUsername(req.GetPhoneNumber())
			if err != nil {
				return nil, status.Errorf(codes.Internal, "failed to found user: %v", err)
			}
		case "email":
			u, err = h.service.GetByUsername(req.GetEmail())
			if err != nil {
				return nil, status.Errorf(codes.Internal, "failed to found user: %v", err)
			}
		case "username":
			u, err = h.service.GetByUsername(req.GetUsername())
			if err != nil {
				return nil, status.Errorf(codes.Internal, "failed to found user: %v", err)
			}
		}
	} else {
		u, err = h.service.GetByID(req.GetId())
	}

	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to fetch user: %v", err)
	}

	if req.GetFirstName() != "" {
		u.FirstName = req.GetFirstName()
	}
	if req.GetLastName() != "" {
		u.LastName = req.GetLastName()
	}
	if req.GetProfileImage() != "" {
		u.ProfileImage = req.GetProfileImage()
	}
	if req.GetGroup() != "" {
		u.Group = req.GetGroup()
	}
	if req.GetTenantId() != "" {
		u.TenantID = req.GetTenantId()
	}

	if req.GetPassword() != "" {
		hashed := util.HashPassword(req.GetPassword())
		u.Password = hashed
	}

	if req.GetId() == "" {
		switch lType {
		case "phone":
			if req.GetUsername() != "" {
				u.Username = req.GetUsername()
			}
			if req.GetEmail() != "" {
				u.Email = req.GetEmail()
			}
		case "email":
			if req.GetUsername() != "" {
				u.Username = req.GetUsername()
			}
			if req.GetPhoneNumber() != "" {
				u.PhoneNumber = req.GetPhoneNumber()
			}
		case "username":
			if req.GetEmail() != "" {
				u.Email = req.GetEmail()
			}
			if req.GetPhoneNumber() != "" {
				u.PhoneNumber = req.GetPhoneNumber()
			}
		case "both":
		}
	} else {
		if req.GetUsername() != "" {
			u.Username = req.GetUsername()
		}
		if req.GetEmail() != "" {
			u.Email = req.GetEmail()
		}
		if req.GetPhoneNumber() != "" {
			u.PhoneNumber = req.GetPhoneNumber()
		}
	}

	if err := h.service.Update(u); err != nil {
		return nil, status.Errorf(codes.Internal, "failed to update user: %v", err)
	}

	return &userpb.UserResponse{
		User: h.toUserPB(u, false),
	}, nil
}

func (h *Handler) ListPaginationUsers(ctx context.Context, req *userpb.ListPaginationUsersRequest) (*userpb.ListPaginationUsersResponse, error) {
	var metaData = cContext.MetaData{
		Limit:       int(req.GetLimit()),
		Sort:        req.GetSort(),
		CurrentPage: int(req.GetPage()),
	}

	metaDataRes, users, err := h.service.PaginationList(metaData, user.UserFilter{SearchFilter: filter.SearchFilter{Search: req.GetSearch()}, TenantID: req.GetTenantId(), IsDelete: "false", Group: req.GetGroup(), NeGroup: req.GetNeGroup()})
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to list users: %v", err)
	}

	var pbusers []*userpb.User
	for _, r := range users {
		userpb := h.toUserPB(&r, req.GetShowID())
		pbusers = append(pbusers, userpb)
	}

	metaDataPB, err := h.toMetaDataPB(&metaDataRes)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to map metadata: %v", err)
	}

	return &userpb.ListPaginationUsersResponse{Users: pbusers, MetaData: metaDataPB}, nil
}

func (h *Handler) BanUser(ctx context.Context, req *userpb.BanUserRequest) (*userpb.UserResponse, error) {

	u, err := h.service.BanUser(req.GetId())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to ban user: %v", err)
	}

	return &userpb.UserResponse{User: h.toUserPB(u, false)}, nil
}

func (h *Handler) GetUser(ctx context.Context, req *userpb.GetUserRequest) (*userpb.UserResponse, error) {
	var (
		u   *user.User
		err error
	)

	switch {
	case req.GetPhoneNumber() != "":
		u, err = h.service.GetByUsername(req.GetPhoneNumber())
	case req.GetEmail() != "":
		u, err = h.service.GetByUsername(req.GetEmail())
	case req.GetUsername() != "":
		u, err = h.service.GetByUsername(req.GetUsername())
	case req.GetId() != "":
		u, err = h.service.GetByID(req.GetId())
	default:
		return nil, status.Errorf(codes.InvalidArgument, "no identifier provided (phone, email, username, or id)")
	}

	if err != nil {
		return nil, status.Errorf(codes.NotFound, "failed to find user: %v", err)
	}

	return &userpb.UserResponse{
		User: h.toUserPB(u, true),
	}, nil
}
