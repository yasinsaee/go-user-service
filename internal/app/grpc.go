package app

import (
	"log"
	"net"

	"github.com/yasinsaee/go-user-service/api/github.com/yasinsaee/go-user-service/api/permissionpb"
	permissiongrpc "github.com/yasinsaee/go-user-service/internal/handlers/grpc/permission"
	repository_permission "github.com/yasinsaee/go-user-service/internal/repository/permission"
	"github.com/yasinsaee/go-user-service/internal/service/permission"
	"github.com/yasinsaee/go-user-service/pkg/mongo"
	"google.golang.org/grpc"
)

func StartGRPCServer() {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()

	permissionRepo := repository_permission.NewMongoPermissionRepository(mongo.DB.Database, "permission")

	permissionService := permission.NewPermissionService(permissionRepo)

	permissionHandler := permissiongrpc.New(permissionService)

	permissionpb.RegisterPermissionServiceServer(s, permissionHandler)

	log.Println("gRPC server is running on port 50051")
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
