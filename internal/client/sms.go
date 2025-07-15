package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/VI-IM/im_backend_go/shared/logger"
)

type SMSClient struct {
	baseURL string
	client  *http.Client
}

type SMSClientInterface interface {
	SendOTP(phone, otp string) error
}

type SMSRequest struct {
	Mobile  string `json:"mobile"`
	Message string `json:"message"`
}

type SMSResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

func NewSMSClient() SMSClientInterface {
	return &SMSClient{
		baseURL: "https://servermsg.com/sendmsg", // This would typically come from config
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

func (s *SMSClient) SendOTP(phone, otp string) error {
	message := fmt.Sprintf("%s is your OTP. Please enter the OTP to verify your mobile number. For more info visit investmango.com", otp)
	
	logger.Get().Info().
		Str("phone", phone).
		Msg("Sending OTP SMS")

	smsRequest := SMSRequest{
		Mobile:  phone,
		Message: message,
	}

	jsonData, err := json.Marshal(smsRequest)
	if err != nil {
		logger.Get().Error().Err(err).Msg("Failed to marshal SMS request")
		return fmt.Errorf("failed to marshal SMS request: %w", err)
	}

	req, err := http.NewRequest("POST", s.baseURL, bytes.NewBuffer(jsonData))
	if err != nil {
		logger.Get().Error().Err(err).Msg("Failed to create SMS request")
		return fmt.Errorf("failed to create SMS request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := s.client.Do(req)
	if err != nil {
		logger.Get().Error().Err(err).Msg("Failed to send SMS")
		return fmt.Errorf("failed to send SMS: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		logger.Get().Error().
			Int("status_code", resp.StatusCode).
			Msg("SMS service returned non-200 status")
		return fmt.Errorf("SMS service returned status: %d", resp.StatusCode)
	}

	var smsResponse SMSResponse
	if err := json.NewDecoder(resp.Body).Decode(&smsResponse); err != nil {
		logger.Get().Error().Err(err).Msg("Failed to decode SMS response")
		return fmt.Errorf("failed to decode SMS response: %w", err)
	}

	logger.Get().Info().
		Str("phone", phone).
		Str("status", smsResponse.Status).
		Msg("SMS sent successfully")

	return nil
}