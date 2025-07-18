package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/VI-IM/im_backend_go/internal/config"
	"github.com/VI-IM/im_backend_go/shared/logger"
)

type CRMClient struct {
	config config.CRM
	client *http.Client
}

type CRMClientInterface interface {
	SendLead(leadData CRMLeadData) error
}

type CRMLeadData struct {
	Name        string `json:"name"`
	Email       string `json:"email"`
	Phone       string `json:"phone"`
	ProjectName string `json:"projectName"`
	QueryInfo   string `json:"queryInfo"`
	Source      string `json:"source"`
}

type CRMResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

func NewCRMClient(cfg config.CRM) CRMClientInterface {
	return &CRMClient{
		config: cfg,
		client: &http.Client{
			Timeout: cfg.Timeout,
		},
	}
}

func (c *CRMClient) SendLead(leadData CRMLeadData) error {
	if !c.config.Enabled {
		logger.Get().Info().Msg("CRM is disabled, skipping lead submission")
		return nil
	}

	logger.Get().Info().
		Str("phone", leadData.Phone).
		Str("email", leadData.Email).
		Msg("Sending lead to external CRM")

	jsonData, err := json.Marshal(leadData)
	if err != nil {
		logger.Get().Error().Err(err).Msg("Failed to marshal CRM lead data")
		return fmt.Errorf("failed to marshal CRM lead data: %w", err)
	}

	var lastErr error
	for attempt := 0; attempt <= c.config.MaxRetries; attempt++ {
		if attempt > 0 {
			logger.Get().Info().
				Int("attempt", attempt).
				Int("max_retries", c.config.MaxRetries).
				Msg("Retrying CRM request")
		}

		req, err := http.NewRequest("POST", c.config.BaseURL, bytes.NewBuffer(jsonData))
		if err != nil {
			logger.Get().Error().Err(err).Msg("Failed to create CRM request")
			return fmt.Errorf("failed to create CRM request: %w", err)
		}

		req.Header.Set("Content-Type", "application/json")

		resp, err := c.client.Do(req)
		if err != nil {
			lastErr = err
			logger.Get().Error().Err(err).
				Int("attempt", attempt+1).
				Msg("Failed to send lead to CRM")
			continue
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
			lastErr = fmt.Errorf("CRM service returned status: %d", resp.StatusCode)
			logger.Get().Error().
				Int("status_code", resp.StatusCode).
				Int("attempt", attempt+1).
				Msg("CRM service returned non-success status")
			continue
		}

		logger.Get().Info().
			Str("phone", leadData.Phone).
			Str("email", leadData.Email).
			Int("attempt", attempt+1).
			Msg("Lead sent to CRM successfully")

		return nil
	}

	logger.Get().Error().Err(lastErr).
		Int("max_retries", c.config.MaxRetries).
		Msg("Failed to send lead to CRM after all retries")
	return fmt.Errorf("failed to send lead to CRM after %d retries: %w", c.config.MaxRetries, lastErr)
}