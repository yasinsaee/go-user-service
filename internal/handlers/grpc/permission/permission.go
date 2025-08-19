package permissiongrpc

import (
	"context"

	"github.com/yasinsaee/go-user-service/internal/domain/permission"
	permissionpb "github.com/yasinsaee/go-user-service/user-service/permission"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Handler struct {
	permissionpb.UnimplementedPermissionServiceServer
	service permission.PermissionService
}

func New(service permission.PermissionService) *Handler {
	return &Handler{service: service}
}

func (h *Handler) CreatePermission(ctx context.Context, req *permissionpb.CreatePermissionRequest) (*permissionpb.CreatePermissionResponse, error) {
	p := &permission.Permission{
		Name:        req.GetName(),
		Description: req.GetDescription(),
	}

	err := h.service.Create(p)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create permission: %v", err)
	}

	return &permissionpb.CreatePermissionResponse{
		Permission: &permissionpb.Permission{
			Id:          p.ID.Hex(),
			Name:        p.Name,
			Description: p.Description,
		},
	}, nil
}

func (h *Handler) UpdatePermission(ctx context.Context, req *permissionpb.UpdatePermissionRequest) (*permissionpb.UpdatePermissionResponse, error) {
	id, err := primitive.ObjectIDFromHex(req.GetId())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid id format")
	}

	per, err := h.service.GetByID(id)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "permission not found: %v", err)
	}

	if name := req.GetName(); name != "" {
		per.Name = name
	}
	if desc := req.GetDescription(); desc != "" {
		per.Description = desc
	}

	if err := h.service.Update(per); err != nil {
		return nil, status.Errorf(codes.Internal, "failed to update permission: %v", err)
	}

	return &permissionpb.UpdatePermissionResponse{
		Permission: &permissionpb.Permission{
			Id:          per.ID.Hex(),
			Name:        per.Name,
			Description: per.Description,
		},
	}, nil
}

func (h *Handler) GetPermission(ctx context.Context, req *permissionpb.GetPermissionRequest) (*permissionpb.GetPermissionResponse, error) {
	id, err := primitive.ObjectIDFromHex(req.GetId())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid id format")
	}

	p, err := h.service.GetByID(id)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "permission not found: %v", err)
	}

	return &permissionpb.GetPermissionResponse{
		Permission: &permissionpb.Permission{
			Id:          p.ID.Hex(),
			Name:        p.Name,
			Description: p.Description,
		},
	}, nil
}

func (h *Handler) ListPermissions(ctx context.Context, req *permissionpb.ListPermissionsRequest) (*permissionpb.ListPermissionsResponse, error) {
	perms, err := h.service.ListAll()
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to list permissions: %v", err)
	}

	var pbPerms []*permissionpb.Permission
	for _, p := range perms {
		pbPerms = append(pbPerms, &permissionpb.Permission{
			Id:          p.ID.Hex(),
			Name:        p.Name,
			Description: p.Description,
		})
	}

	return &permissionpb.ListPermissionsResponse{
		Permissions: pbPerms,
	}, nil
}

func (h *Handler) DeletePermission(ctx context.Context, req *permissionpb.DeletePermissionRequest) (*permissionpb.DeletePermissionResponse, error) {
	err := h.service.Delete(req.GetId())
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "permission not found: %v", err)
	}

	return &permissionpb.DeletePermissionResponse{
		Message: "ok",
	}, nil
}
