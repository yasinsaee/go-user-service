package role

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Role struct {
	ID          primitive.ObjectID `bson:"_id,omitempty"`
	UniqueID    string             `bson:"unique_id"`
	Scope       string             `bson:"scope"` //can be null - Logical separation (e.g., 'global' or 'project')
	Name        string             `bson:"name"`
	Key         string             `bson:"key"`
	Description string             `bson:"description,omitempty"`
	Permissions []string           `bson:"permissions"`
	IsDeleted   bool               `bson:"is_deleted"`
	DeletedAt   time.Time          `bson:"deleted_at"`
	CreatedAt   time.Time          `bson:"created_at"`
	UpdatedAt   time.Time          `bson:"updated_at"`
}

type Roles []Role
