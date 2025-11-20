package otp

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type (
	Otp struct {
		ID          primitive.ObjectID `bson:"_id,omitempty" json:"id"`
		Receiver    string             `bson:"receiver" json:"receiver"`
		Code        string             `bson:"code" json:"code"`
		Description string             `bson:"description" json:"description"`
		Used        bool               `bson:"used" json:"used"`
		SendAt      time.Time          `bson:"send_at" json:"send_at"`
		ExpiresAt   time.Time          `bson:"expires_at" json:"expires_at"`
		CreatedAt   time.Time          `bson:"created_at" json:"created_at"`
		UpdatedAt   time.Time          `bson:"updated_at" json:"updated_at"`
	}
)

type Otps []Otp
