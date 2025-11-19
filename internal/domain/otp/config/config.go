package config

import (
	"os"
	"strconv"
	"time"
)

type OTPConfig struct {
	Length    int
	Charset   string
	TTL       time.Duration
	RateLimit time.Duration
}

func LoadOTPConfig() OTPConfig {
	length := getEnvInt("OTP_LENGTH", 6)
	charset := getEnv("OTP_CHARSET", "0123456789")
	ttl := time.Duration(getEnvInt("OTP_TTL_SECONDS", 120)) * time.Second
	rateLimit := time.Duration(getEnvInt("OTP_RATE_LIMIT_SECONDS", 60)) * time.Second

	return OTPConfig{
		Length:    length,
		Charset:   charset,
		TTL:       ttl,
		RateLimit: rateLimit,
	}
}

func getEnv(key, defaultVal string) string {
	if val, ok := os.LookupEnv(key); ok {
		return val
	}
	return defaultVal
}

func getEnvInt(key string, defaultVal int) int {
	if val, ok := os.LookupEnv(key); ok {
		if i, err := strconv.Atoi(val); err == nil {
			return i
		}
	}
	return defaultVal
}
