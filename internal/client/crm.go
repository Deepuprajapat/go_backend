package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/VI-IM/im_backend_go/shared/logger"
)

type CRMClient struct {
	baseURL string
	client  *http.Client
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

func NewCRMClient() CRMClientInterface {
	return &CRMClient{
		baseURL: "http://148.66.133.154:8181/new-leads/from/open-source",
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

func (c *CRMClient) SendLead(leadData CRMLeadData) error {
	logger.Get().Info().
		Str("phone", leadData.Phone).
		Str("email", leadData.Email).
		Msg("Sending lead to external CRM")

	jsonData, err := json.Marshal(leadData)
	if err != nil {
		logger.Get().Error().Err(err).Msg("Failed to marshal CRM lead data")
		return fmt.Errorf("failed to marshal CRM lead data: %w", err)
	}

	req, err := http.NewRequest("POST", c.baseURL, bytes.NewBuffer(jsonData))
	if err != nil {
		logger.Get().Error().Err(err).Msg("Failed to create CRM request")
		return fmt.Errorf("failed to create CRM request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := c.client.Do(req)
	if err != nil {
		logger.Get().Error().Err(err).Msg("Failed to send lead to CRM")
		return fmt.Errorf("failed to send lead to CRM: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		logger.Get().Error().
			Int("status_code", resp.StatusCode).
			Msg("CRM service returned non-success status")
		return fmt.Errorf("CRM service returned status: %d", resp.StatusCode)
	}

	logger.Get().Info().
		Str("phone", leadData.Phone).
		Str("email", leadData.Email).
		Msg("Lead sent to CRM successfully")

	return nil
}