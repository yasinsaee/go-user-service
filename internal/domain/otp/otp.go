package otp

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type (
	OTP struct {
		ID          primitive.ObjectID `bson:"_id,omitempty" json:"id"`
		UserID      primitive.ObjectID `bson:"user_id" json:"user_id"`
		Code        string             `bson:"code" json:"code"`
		Description string             `bson:"description" json:"description"`
		Used        bool               `bson:"used" json:"used"`
		SendAt      time.Time          `bson:"send_at" json:"send_at"`
		ExpiresAt   time.Time          `bson:"expires_at" json:"expires_at"`
		CreatedAt   time.Time          `bson:"created_at" json:"created_at"`
		UpdatedAt   time.Time          `bson:"updated_at" json:"updated_at"`
	}
)

type OTPs []OTP
