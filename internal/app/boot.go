package app

import (
	"log"
	"strings"

	"github.com/yasinsaee/go-user-service/internal/domain/group"
	"github.com/yasinsaee/go-user-service/internal/domain/permission"
	"github.com/yasinsaee/go-user-service/internal/domain/role"
	"github.com/yasinsaee/go-user-service/internal/domain/user"
)

func BootstrapSystem(
	userService user.UserService,
	roleService role.RoleService,
	permService permission.PermissionService,
	groupService group.GroupService,
	adminEmail, adminPass string,
) {
	log.Println("🌱 Bootstrapping system data...")

	permissionsDef := map[string]string{
		"create": "Create new resources",
		"read":   "View and read resources",
		"update": "Edit existing resources",
		"delete": "Remove resources",
		"import": "Import data from external sources",
		"export": "Export data to external files",
	}

	permIDs := make([]string, 0, len(permissionsDef))

	for key, desc := range permissionsDef {
		p, _ := permService.GetByKey(key)

		if p == nil {
			newPerm := &permission.Permission{
				Key:         key,
				Name:        strings.Title(key),
				Description: desc,
			}
			err := permService.Create(newPerm)
			if err != nil {
				log.Printf("❌ Failed to create permission %s: %v", key, err)
				continue
			}
			permIDs = append(permIDs, newPerm.UniqueID)
		} else {
			permIDs = append(permIDs, p.UniqueID)
		}
	}

	superAdminRole, _ := roleService.GetByKey("super_admin")
	if superAdminRole == nil {
		superAdminRole = &role.Role{
			Key:         "super_admin",
			Name:        "Super Admin",
			Description: "Full access to all system resources",
			Permissions: permIDs,
		}
		err := roleService.Create(superAdminRole)
		if err != nil {
			log.Fatalf("❌ Failed to create super admin role: %v", err)
		}
	}

	adminGroup, _ := groupService.GetByName("Administrators")
	if adminGroup == nil {
		adminGroup = &group.Group{
			Name:        "Administrators",
			Description: "System administrators group",
			SelectedRoles: []group.SelectedRole{
				{
					RoleID:            superAdminRole.UniqueID,
					DeniedPermissions: []string{},
				},
			},
		}
		err := groupService.Create(adminGroup)
		if err != nil {
			log.Fatalf("❌ Failed to create admin group: %v", err)
		}
	}

	existingUser, _ := userService.GetByUsername(adminEmail)
	if existingUser == nil {
		adminUser := &user.User{
			Password:  adminPass,
			FirstName: "Super",
			LastName:  "Admin",
			IsActive:  true,
			Group:     adminGroup.UniqueID,
			Username:  adminEmail,
			Email:     "admin@admin.com",
		}

		err := userService.Register(adminEmail, adminUser)
		if err != nil {
			log.Fatalf("❌ Failed to create admin user: %v", err)
		}

		log.Println("✅ Super Admin user created successfully.")
	} else {
		log.Println("ℹ️ Super Admin already exists.")
	}
}
