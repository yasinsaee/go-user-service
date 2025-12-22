package otp

type OtpType int32

const (
	OtpTypeRegister OtpType = iota
	OtpTypeForgotPassword
	OtpTypeVerifyEmail
)