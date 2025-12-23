package otp

import "go.mongodb.org/mongo-driver/bson"

type OTPRepository interface {
	Create(otp *Otp) error
	FindByID(id any) (*Otp, error)
	FindByName(name string) (*Otp, error)
	Update(otp *Otp) error
	Delete(id any) error
	List() (Otps, error)
	FindByReceiverAndCode(receiver, code string) (*Otp, error)
	FindByReceiverAndCodeAndType(receiver, code string, otpType OtpType) (*Otp, error)
	ListByType(otpType OtpType) (Otps, error)
	DeleteExpiredOtps() error
	Count(q bson.M) (int, error)
}
