package user

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID           primitive.ObjectID   `bson:"_id,omitempty" json:"id"`
	Roles        []primitive.ObjectID `bson:"roles" json:"roles"`
	FirstName    string               `bson:"first_name" json:"first_name"`
	LastName     string               `bson:"last_name" json:"last_name"`
	ProfileImage string               `bson:"profile_image,omitempty" json:"profile_image,omitempty"`
	Username     string               `bson:"username" json:"username"`
	Email        string               `bson:"email" json:"email"`
	Password     string               `bson:"password" json:"-"`
	PhoneNumber  string               `bson:"phone_number,omitempty" json:"phone_number,omitempty"` //optional
	IsActive     bool                 `bson:"is_active" json:"is_active"`
	CreatedAt    time.Time            `bson:"created_at" json:"created_at"`
	UpdatedAt    time.Time            `bson:"updated_at" json:"updated_at"`
	LastLogin    *time.Time           `bson:"last_login,omitempty" json:"last_login,omitempty"`
}

type Users []User
