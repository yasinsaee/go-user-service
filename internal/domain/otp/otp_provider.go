package otp

type OTPProvider interface {
	Send(receiver string, code string) error
}
