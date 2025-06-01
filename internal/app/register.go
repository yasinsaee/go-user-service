package app

import (
	"github.com/labstack/echo/v4"
	handler_permission "github.com/yasinsaee/go-user-service/internal/handlers/rest/permission"
	role_permission "github.com/yasinsaee/go-user-service/internal/handlers/rest/role"
	user_permission "github.com/yasinsaee/go-user-service/internal/handlers/rest/user"
	repository_permission "github.com/yasinsaee/go-user-service/internal/repository/permission"
	repository_role "github.com/yasinsaee/go-user-service/internal/repository/role"
	repository_user "github.com/yasinsaee/go-user-service/internal/repository/user"
	"github.com/yasinsaee/go-user-service/internal/service/permission"
	"github.com/yasinsaee/go-user-service/internal/service/role"
	"github.com/yasinsaee/go-user-service/internal/service/user"
	"github.com/yasinsaee/go-user-service/pkg/mongo"
)

func Register(e *echo.Echo) {
	permissionRepo := repository_permission.NewMongoPermissionRepository(mongo.DB.Database, "permission")
	roleRepo := repository_role.NewMongoRoleRepository(mongo.DB.Database, "role")
	userRepo := repository_user.NewMongoUserRepository(mongo.DB.Database, "user")

	permissionService := permission.NewPermissionService(permissionRepo)
	roleService := role.NewRoleService(roleRepo)
	userService := user.NewUserService(userRepo)

	permissionHandler := handler_permission.NewPermissionHandler(permissionService)
	roleHandler := role_permission.NewRoleHandler(roleService)
	userHandler := user_permission.NewUserHandler(userService)

	permissionHandler.RegisterRoutes(e)
	roleHandler.RegisterRoutes(e)
	userHandler.RegisterRoutes(e)
}
