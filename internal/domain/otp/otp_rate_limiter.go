package otp

import "time"

type OTPRateLimiter interface {
	CanSend(receiver string) (bool, error)
	MarkSend(receiver string, ttl time.Duration) error
}