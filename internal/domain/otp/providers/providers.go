package providers

import (
	"strconv"

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
		// SMTP provider

	case "SMSIR":
		apiKey := config.GetEnv("SMSIR_API_KEY", "")
		senderStr := config.GetEnv("SMSIR_SENDER", "")
		message := config.GetEnv(
			"SMSIR_MESSAGE",
			"به یوزر سرویس خوش اومدی! 🌿\nخوشحالیم که کنار مایی.\nکد تایید ورود شما: {{code}}\nاین کد رو در اختیار دیگران قرار نده.",
		)

		sender, err := strconv.Atoi(senderStr)
		if err != nil {
			return nil
		}

		return NewSMSProviderSmsIR(apiKey, sender, message)
	}

	return nil
}
