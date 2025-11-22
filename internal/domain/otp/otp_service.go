package otp

type OTPService interface {

	// CRUD
	Create(otp *Otp) error
	GetByID(id any) (*Otp, error)
	GetByName(name string) (*Otp, error)
	Update(otp *Otp) error
	Delete(id any) error
	ListAll() (Otps, error)

	// OTP business logic
	GenerateCode() string
	SaveCode(receiver string, code string) error
	ValidateCode(receiver string, code string) (bool, error)
	SendCode(receiver string, code string) error

	// Rate limiting
	CanSend(receiver string) (bool, error)
	MarkSend(receiver string) error
}
