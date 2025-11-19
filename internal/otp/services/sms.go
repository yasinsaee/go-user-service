package services

//Sms.Ir

type SmsOTP struct {
	APIKey string
}

func (k *SmsOTP) SendOTP(to, code string) error {
	// Call Sms API
	return nil
}
