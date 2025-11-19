package otp

import "fmt"

type MultiOTP struct {
	Providers []OTPService
}

func (m *MultiOTP) SendOTP(to, code string) error {
	for _, p := range m.Providers {
		if err := p.SendOTP(to, code); err == nil {
			return nil
		}
	}
	return fmt.Errorf("all otp providers failed")
}
