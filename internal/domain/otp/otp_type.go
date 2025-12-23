package otp

type OtpType int32

const (
	OtpTypeForgotPassword OtpType = iota
	OtpTypeRegister
	OtpTypeVerifyEmail
)