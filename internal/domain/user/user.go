package user

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID           primitive.ObjectID `bson:"_id,omitempty"`
	UniqueID     string             `bson:"unique_id"`
	TenantID     string             `bson:"tenant_id"` //can be null - Organization or Store ID for multi-tenancy
	FirstName    string             `bson:"first_name"`
	LastName     string             `bson:"last_name"`
	ProfileImage string             `bson:"profile_image,omitempty"`
	Username     string             `bson:"username"`
	Email        string             `bson:"email"`
	Password     string             `bson:"password"`
	PhoneNumber  string             `bson:"phone_number"`
	Group        string             `bson:"group"`
	IsActive     bool               `bson:"is_active"`
	IsBanned     bool               `bson:"is_banned"`
	IsDeleted    bool               `bson:"is_deleted"`
	BannedAt     time.Time          `bson:"banned_at"`
	DeletedAt    time.Time          `bson:"deleted_at"`
	CreatedAt    time.Time          `bson:"created_at"`
	UpdatedAt    time.Time          `bson:"updated_at"`
	LastLogin    time.Time          `bson:"last_login,omitempty"`
}

type Users []User
