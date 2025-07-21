package response

import (
	"time"

	"github.com/VI-IM/im_backend_go/ent"
)

type Lead struct {
	ID                   int       `json:"id"`
	Name                 string    `json:"name"`
	Email                string    `json:"email"`
	Phone                string    `json:"phone"`
	Message              string    `json:"message,omitempty"`
	Source               string    `json:"source"`
	IsDuplicate          bool      `json:"is_duplicate"`
	DuplicateReferenceID string    `json:"duplicate_reference_id,omitempty"`
	OtpVerified          bool      `json:"otp_verified"`
	SyncStatus           string    `json:"sync_status"`
	PropertyID           string    `json:"property_id,omitempty"`
	ProjectID            string    `json:"project_id,omitempty"`
	ProjectName          string    `json:"project_name,omitempty"`
	PropertyName         string    `json:"property_name,omitempty"`
	CreatedAt            time.Time `json:"created_at"`
	UpdatedAt            time.Time `json:"updated_at"`
}

type LeadListResponse struct {
	Content []*Lead `json:"content"`
}

type CreateLeadResponse struct {
	Message string `json:"message"`
}

type ValidateOTPResponse struct {
	Message string `json:"message"`
}

type ResendOTPResponse struct {
	Message string `json:"message"`
}

type DuplicateLeadGroup struct {
	Last    *Lead   `json:"last"`
	History []*Lead `json:"history"`
}

type DateLeadsResponse struct {
	Date map[string]*DateLeadsData `json:"date"`
}

type DateLeadsData struct {
	UniqueLeads    []*Lead                       `json:"unique_leads"`
	DuplicateLeads map[string]*DuplicateLeadGroup `json:"duplicate_leads"`
}

func ToLeadResponse(lead *ent.Leads) *Lead {
	response := &Lead{
		ID:                   lead.ID,
		Name:                 lead.Name,
		Email:                lead.Email,
		Phone:                lead.Phone,
		Message:              lead.Message,
		Source:               lead.Source,
		IsDuplicate:          lead.IsDuplicate,
		DuplicateReferenceID: lead.DuplicateReferenceID,
		OtpVerified:          lead.OtpVerified,
		SyncStatus:           string(lead.SyncStatus),
		CreatedAt:            lead.CreatedAt,
		UpdatedAt:            lead.UpdatedAt,
	}

	if lead.Edges.Property != nil {
		response.PropertyID = lead.Edges.Property.ID
		response.PropertyName = lead.Edges.Property.Name
	}

	if lead.Edges.Project != nil {
		response.ProjectID = lead.Edges.Project.ID
		response.ProjectName = lead.Edges.Project.Name
	}

	// Get project name from property->project relationship if direct project is not available
	if response.ProjectName == "" && lead.Edges.Property != nil && lead.Edges.Property.Edges.Project != nil {
		response.ProjectName = lead.Edges.Property.Edges.Project.Name
	}

	return response
}

func ToLeadListResponse(leads []*ent.Leads) *LeadListResponse {
	content := make([]*Lead, len(leads))
	for i, lead := range leads {
		content[i] = ToLeadResponse(lead)
	}

	return &LeadListResponse{
		Content: content,
	}
}