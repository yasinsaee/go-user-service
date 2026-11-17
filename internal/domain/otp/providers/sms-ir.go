package providers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

type SMSProviderSmsIR struct {
	APIKey     string
	LineNumber int
	Message    string
	HTTPClient *http.Client
}

func NewSMSProviderSmsIR(apiKey string, sender int, message string) *SMSProviderSmsIR {
	return &SMSProviderSmsIR{
		APIKey:     apiKey,
		LineNumber: sender,
		Message:    message,
		HTTPClient: http.DefaultClient,
	}
}

func (s *SMSProviderSmsIR) Send(receiver string, code string) error {
	const apiURL = "https://api.sms.ir/v1/send/bulk"

	message := strings.ReplaceAll(s.Message, "{{code}}", code)

	payload := struct {
		LineNumber  int      `json:"lineNumber"`
		MessageText string   `json:"messageText"`
		Mobiles     []string `json:"mobiles"`
	}{
		LineNumber:  s.LineNumber,
		MessageText: message,
		Mobiles:     []string{receiver},
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("marshal sms payload: %w", err)
	}

	req, err := http.NewRequest(http.MethodPost, apiURL, bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("create sms request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("X-API-KEY", s.APIKey)

	resp, err := s.HTTPClient.Do(req)
	if err != nil {
		return fmt.Errorf("send sms request: %w", err)
	}
	defer resp.Body.Close()

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("read sms response: %w", err)
	}

	fmt.Println("SMS.ir Status:", resp.StatusCode)
	fmt.Println("SMS.ir Response:", string(responseBody))

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf(
			"sms sending failed: status %d, response: %s",
			resp.StatusCode,
			string(responseBody),
		)
	}

	return nil
}