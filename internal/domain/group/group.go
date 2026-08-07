package group

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type (
	Group struct {
		ID            primitive.ObjectID `bson:"_id,omitempty"`
		Owner         string             `bson:"owner"`
		UniqueID      string             `bson:"unique_id"`
		TenantID      string             `bson:"tenant_id"` //can be null - Organization or Store ID for multi-tenancy
		Name          string             `bson:"name"`
		Key           string             `bson:"key"`
		SelectedRoles SelectedRoles      `bson:"selected_roles"`
		Description   string             `bson:"description,omitempty"`
		IsDeleted     bool               `bson:"is_deleted"`
		DeletedAt     time.Time          `bson:"deleted_at"`
		CreatedAt     time.Time          `bson:"created_at"`
		UpdatedAt     time.Time          `bson:"updated_at"`
	}

	SelectedRole struct {
		RoleID            string   `bson:"role_id"`                      // Reference to the Role ID
		DeniedPermissions []string `bson:"denied_permissions,omitempty"` // List of permissions explicitly denied
	}
)
type SelectedRoles []SelectedRole
type Groups []Group
