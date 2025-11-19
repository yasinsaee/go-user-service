package providers

import (
	"fmt"
	"net/http"
	"net/url"
)

type SMSProviderKavenegar struct {
	APIKey     string
	Sender     string
	HTTPClient *http.Client
}

func NewSMSProviderKavenegar(apiKey, sender string) *SMSProviderKavenegar {
	return &SMSProviderKavenegar{
		APIKey:     apiKey,
		Sender:     sender,
		HTTPClient: http.DefaultClient,
	}
}


//Should Change with Document Im not sure this is work
//please if you want another provider send document or send pr
//yasinvsaee@gmail.com
func (s *SMSProviderKavenegar) Send(receiver string, code string) error {
	message := fmt.Sprintf("Your OTP code is: %s", code)

	apiURL := fmt.Sprintf("https://api.kavenegar.com/v1/%s/sms/send.json", s.APIKey)

	data := url.Values{}
	data.Set("receptor", receiver)
	data.Set("message", message)
	data.Set("sender", s.Sender)

	resp, err := s.HTTPClient.PostForm(apiURL, data)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return fmt.Errorf("sms sending failed: status %d", resp.StatusCode)
	}

	return nil
}
