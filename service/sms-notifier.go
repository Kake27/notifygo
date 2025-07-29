package service
import (
	"fmt"
	"net/http"
	"net/url"
	"io"
	"notification-service/config"
)


type MockSMSService struct {}
func NewMockSMSService() *MockSMSService {
	return &MockSMSService{}
}

func (m *MockSMSService) Send(to, subject, message string) (string, error) {
	fmt.Printf("[MOCK SMS] To: %s | Subject: %s | Message: %s\n", to, subject, message)
    return fmt.Sprintf("Mock SMS sent to %s", to), nil
}

type SMSService struct {

	cfg config.SMSConfig
}

func NewSMSService() *SMSService {
	return &SMSService{
		cfg : config.GetSMSConfig(),
	}
}

func (s *SMSService) Send(to, message string) (string, error) {
	params := url.Values{}
	params.Set("username", s.cfg.Username)
    params.Set("password", s.cfg.Password)
    params.Set("to", to)
	params.Set("from", s.cfg.Sender)
    params.Set("text", message)

	finalURL := fmt.Sprintf("%s?%s", s.cfg.GatewayURL, params.Encode())
	resp, err := http.Get(finalURL)
    if err != nil {
        return "", fmt.Errorf("failed to send SMS: %w", err)
    }
    defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
    if resp.StatusCode<200 || resp.StatusCode>299 {
        return "", fmt.Errorf("SMS gateway returned %d: %s", resp.StatusCode, string(body))
    }

    return fmt.Sprintf("SMS successfully sent to %s", to), nil
}