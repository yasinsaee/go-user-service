package permission

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Permission struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	UniqueID    string             `bson:"unique_id" json:"unique_id"`
	Key         string             `bson:"key" json:"key"`
	Name        string             `bson:"name" json:"name"`
	Description string             `bson:"description,omitempty" json:"description,omitempty"`
	IsDeleted   bool               `bson:"is_deleted"`
	DeletedAt   time.Time          `bson:"deleted_at"`
	CreatedAt   time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt   time.Time          `bson:"updated_at" json:"updated_at"`
}

type Permissions []Permission
