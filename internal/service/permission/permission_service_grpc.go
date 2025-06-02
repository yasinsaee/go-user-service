package permission

import (
	"github.com/yasinsaee/go-user-service/api/github.com/yasinsaee/go-user-service/api/permissionpb"
	"github.com/yasinsaee/go-user-service/internal/domain/permission"
)

type PermissionGRPCServer struct {
	permissionpb.UnimplementedPermissionServiceServer
	service permission.PermissionService
}

	