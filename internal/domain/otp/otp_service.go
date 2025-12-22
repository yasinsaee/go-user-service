package otp

import "go.mongodb.org/mongo-driver/bson"

type OTPService interface {

	// CRUD
	Create(otp *Otp) error
	GetByID(id any) (*Otp, error)
	GetByName(name string) (*Otp, error)
	Update(otp *Otp) error
	Delete(id any) error
	ListAll() (Otps, error)
	ListByType(otpType OtpType) (Otps, error)
	Count(q bson.M) (int, error)

	// OTP business logic
	GenerateCode() string
	SaveCode(receiver string, otpType OtpType, code string) error
	ValidateCode(receiver string, otpType OtpType, code string) (bool, error)
	SendCode(receiver string, otpType OtpType, code string) error

	// Rate limiting
	CanSend(receiver string) (bool, error)
	MarkSend(receiver string) error

	// Hard limit Check
	CheckHardLimit(receiver string, otpType OtpType) (bool, error)
}
