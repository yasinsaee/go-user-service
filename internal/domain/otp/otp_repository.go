package otp

import "go.mongodb.org/mongo-driver/bson"

type OTPRepository interface {
	Create(OTP *Otp) error
	FindByID(id any) (*Otp, error)
	FindByName(name string) (*Otp, error)
	Update(OTP *Otp) error
	Delete(id any) error
	List() (Otps, error)
	FindByReceiverAndCode(receiver, code string) (*Otp, error)
	DeleteExpiredOtps() error
	Count(q bson.M) (int, error)
}
