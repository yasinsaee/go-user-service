package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Role struct {
	ID          primitive.ObjectID   `bson:"_id,omitempty" json:"id"`
	Permissions []primitive.ObjectID `bson:"permissions" json:"permissions"`
	Name        string               `bson:"name" json:"name"`
	Description string               `bson:"description,omitempty" json:"description,omitempty"`
}
