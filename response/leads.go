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
	PropertyID           string    `json:"property_id,omitempty"`
	ProjectID            string    `json:"project_id,omitempty"`
	CreatedAt            time.Time `json:"created_at"`
	UpdatedAt            time.Time `json:"updated_at"`
}

type LeadListResponse struct {
	Content      []*Lead `json:"content"`
	TotalElements int     `json:"total_elements"`
	TotalPages    int     `json:"total_pages"`
	Size          int     `json:"size"`
	Number        int     `json:"number"`
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
		CreatedAt:            lead.CreatedAt,
		UpdatedAt:            lead.UpdatedAt,
	}

	if lead.Edges.Property != nil {
		response.PropertyID = lead.Edges.Property.ID
	}

	if lead.Edges.Project != nil {
		response.ProjectID = lead.Edges.Project.ID
	}

	return response
}

func ToLeadListResponse(leads []*ent.Leads, totalElements, size, page int) *LeadListResponse {
	content := make([]*Lead, len(leads))
	for i, lead := range leads {
		content[i] = ToLeadResponse(lead)
	}

	totalPages := (totalElements + size - 1) / size

	return &LeadListResponse{
		Content:       content,
		TotalElements: totalElements,
		TotalPages:    totalPages,
		Size:          size,
		Number:        page,
	}
}