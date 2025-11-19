package otp

type OTPService interface {

	// CRUD
	Create(otp *OTP) error
	GetByID(id any) (*OTP, error)
	GetByName(name string) (*OTP, error)
	Update(otp *OTP) error
	Delete(id any) error
	ListAll() (OTPs, error)

	// OTP business logic
	GenerateCode() string
	SaveCode(receiver string, code string, ttlSeconds int) error
	ValidateCode(receiver string, code string) (bool, error)
	SendCode(receiver string, code string) error

	// Rate limiting
	CanSend(receiver string) (bool, error)
	MarkSend(receiver string) error
}
