package otp

type OTPRepository interface {
	Create(OTP *Otp) error
	FindByID(id any) (*Otp, error)
	FindByName(name string) (*Otp, error)
	Update(OTP *Otp) error
	Delete(id any) error
	List() (Otps, error)
}
