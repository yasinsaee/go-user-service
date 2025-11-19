package otp

type OTPRepository interface {
	Create(OTP *OTP) error
	FindByID(id any) (*OTP, error)
	FindByName(name string) (*OTP, error)
	Update(OTP *OTP) error
	Delete(id any) error
	List() (OTPs, error)
}
