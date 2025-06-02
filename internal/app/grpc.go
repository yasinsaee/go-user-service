package app

import (
	"log"
	"net"

	"github.com/yasinsaee/go-user-service/api/permissionpb"
	"github.com/yasinsaee/go-user-service/api/rolepb"
	permissiongrpc "github.com/yasinsaee/go-user-service/internal/handlers/grpc/permission"
	rolegrpc "github.com/yasinsaee/go-user-service/internal/handlers/grpc/role"
	repository_permission "github.com/yasinsaee/go-user-service/internal/repository/permission"
	repository_role "github.com/yasinsaee/go-user-service/internal/repository/role"
	"github.com/yasinsaee/go-user-service/internal/service/permission"
	"github.com/yasinsaee/go-user-service/internal/service/role"
	"github.com/yasinsaee/go-user-service/pkg/mongo"
	"google.golang.org/grpc"
)

func StartGRPCServer() {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()

	//repos
	permissionRepo := repository_permission.NewMongoPermissionRepository(mongo.DB.Database, "permission")
	roleRepo := repository_role.NewMongoRoleRepository(mongo.DB.Database, "role")

	//services
	permissionService := permission.NewPermissionService(permissionRepo)
	roleService := role.NewRoleService(roleRepo)

	//handlers
	permissionHandler := permissiongrpc.New(permissionService)
	roleHandler := rolegrpc.New(roleService)

	//register grpc services
	permissionpb.RegisterPermissionServiceServer(s, permissionHandler)
	rolepb.RegisterRoleServiceServer(s, roleHandler )

	log.Println("gRPC server is running on port 50051")
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
