package otp

type OTPService interface {
	SendOTP(to string, code string) error
}
