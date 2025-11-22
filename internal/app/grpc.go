package app

import (
	"log"
	"net"

	"github.com/yasinsaee/go-user-service/internal/app/config"
	otp_config "github.com/yasinsaee/go-user-service/internal/domain/otp/config"
	"github.com/yasinsaee/go-user-service/internal/domain/otp/providers"
	otpgrpc "github.com/yasinsaee/go-user-service/internal/handlers/grpc/otp"
	permissiongrpc "github.com/yasinsaee/go-user-service/internal/handlers/grpc/permission"
	rolegrpc "github.com/yasinsaee/go-user-service/internal/handlers/grpc/role"
	usergrpc "github.com/yasinsaee/go-user-service/internal/handlers/grpc/user"
	repository_otp "github.com/yasinsaee/go-user-service/internal/repository/otp"
	repository_permission "github.com/yasinsaee/go-user-service/internal/repository/permission"
	repository_role "github.com/yasinsaee/go-user-service/internal/repository/role"
	repository_user "github.com/yasinsaee/go-user-service/internal/repository/user"
	"github.com/yasinsaee/go-user-service/internal/service/otp"
	ratelimiter "github.com/yasinsaee/go-user-service/internal/service/otp/redis"
	"github.com/yasinsaee/go-user-service/internal/service/permission"
	"github.com/yasinsaee/go-user-service/internal/service/role"
	"github.com/yasinsaee/go-user-service/internal/service/user"
	"github.com/yasinsaee/go-user-service/pkg/mongo"
	otppb "github.com/yasinsaee/go-user-service/user-service/otp"
	permissionpb "github.com/yasinsaee/go-user-service/user-service/permission"
	rolepb "github.com/yasinsaee/go-user-service/user-service/role"
	userpb "github.com/yasinsaee/go-user-service/user-service/user"
	"google.golang.org/grpc"
)

func StartGRPCServer() {
	port := config.GetEnv("PORT", "50051")
	lis, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	// s := grpc.NewServer(grpc.UnaryInterceptor(middleware.AuthInterceptor()))
	s := grpc.NewServer()

	//repos
	permissionRepo := repository_permission.NewMongoPermissionRepository(mongo.DB.Database, "permission")
	roleRepo := repository_role.NewMongoRoleRepository(mongo.DB.Database, "role")
	userRepo := repository_user.NewMongoUserRepository(mongo.DB.Database, "user")
	otpRepo := repository_otp.NewMongoOTPRepository(mongo.DB.Database, "otp")

	//providers
	provider := providers.NewOTPProvider()

	//otp config
	otpConfig := otp_config.LoadOTPConfig()
	////rate limiter
	rateLimiter := ratelimiter.NewRedisOTPRateLimiter(int(otpConfig.RateLimit))

	//services
	permissionService := permission.NewPermissionService(permissionRepo)
	roleService := role.NewRoleService(roleRepo)
	userService := user.NewUserService(userRepo)
	otpService := otp.NewOTPService(otpRepo, provider, rateLimiter, otpConfig.TTL, otpConfig.RateLimit, otpConfig, otpConfig.MaxOTPPerReceiver)

	//handlers
	permissionHandler := permissiongrpc.New(permissionService)
	roleHandler := rolegrpc.New(roleService, permissionService)
	userHandler := usergrpc.New(userService, roleService, permissionService)
	otpHandler := otpgrpc.New(otpService)

	//register grpc services
	permissionpb.RegisterPermissionServiceServer(s, permissionHandler)
	rolepb.RegisterRoleServiceServer(s, roleHandler)
	userpb.RegisterUserServiceServer(s, userHandler)
	otppb.RegisterOTPServiceServer(s, otpHandler)

	log.Println("gRPC server is running on port 50051")
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
