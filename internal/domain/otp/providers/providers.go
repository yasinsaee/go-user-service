package providers

import (
	"github.com/yasinsaee/go-user-service/internal/app/config"
	"github.com/yasinsaee/go-user-service/internal/domain/otp"
)

func NewOTPProvider() otp.OTPProvider {
	providerType := config.GetEnv("OTP_PROVIDER", "")

	switch providerType {
	case "KAVENEGAR":
		apiKey := config.GetEnv("KAVENEGAR_API_KEY", "")
		sender := config.GetEnv("KAVENEGAR_SENDER", "")
		return NewSMSProviderKavenegar(apiKey, sender)

	case "SMTP":
		// host := config.GetEnv("SMTP_HOST", "")
		// port := config.GetEnv("SMTP_PORT", "587")
		// user := config.GetEnv("SMTP_USER", "")
		// pass := config.GetEnv("SMTP_PASS", "")
		// return NewEmailProvider(host, port, user, pass)
	}
	return nil
}
