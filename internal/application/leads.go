package application

import (
	"context"
	"fmt"
	"math/rand"
	"net/http"
	"strconv"
	"time"

	"github.com/VI-IM/im_backend_go/ent"
	"github.com/VI-IM/im_backend_go/ent/leads"
	"github.com/VI-IM/im_backend_go/internal/client"
	"github.com/VI-IM/im_backend_go/request"
	"github.com/VI-IM/im_backend_go/response"
	imhttp "github.com/VI-IM/im_backend_go/shared"
	"github.com/VI-IM/im_backend_go/shared/logger"
)

func generateOTP() string {
	rand.Seed(time.Now().UnixNano())
	return fmt.Sprintf("%06d", rand.Intn(1000000))
}

func (a *application) CreateLeadWithOTP(ctx context.Context, req *request.CreateLeadRequest) (*response.CreateLeadResponse, *imhttp.CustomError) {

	// Check for existing lead by phone number
	existingLead, err := a.repo.GetLeadByPhone(ctx, req.Phone)

	lead := &ent.Leads{
		Email:       req.Email,
		Name:        req.Name,
		Phone:       req.Phone,
		Message:     req.Message,
		Source:      "Organic",
		IsDuplicate: false,
		OtpVerified: false,
	}

	// Handle duplicate detection
	if err == nil && existingLead != nil {
		logger.Get().Info().Str("phone", req.Phone).Msg("Duplicate lead detected")
		lead.IsDuplicate = true
		lead.DuplicateReferenceID = strconv.Itoa(existingLead.ID)
	}

	// Set property or project relationship
	if req.PropertyID != "" {
		property, err := a.repo.GetPropertyByID(req.PropertyID)
		if err != nil {
			logger.Get().Warn().Err(err).Str("property_id", req.PropertyID).Msg("Property not found, creating lead without property association")
		} else {
			lead.Edges.Property = property
		}
	}

	if req.ProjectID != "" {
		project, err := a.repo.GetProjectByID(req.ProjectID)
		if err != nil {
			logger.Get().Warn().Err(err).Str("project_id", req.ProjectID).Msg("Project not found, creating lead without project association")
		} else {
			lead.Edges.Project = project
		}
	}

	// Generate and set OTP
	otp := generateOTP()
	lead.Otp = otp

	// Save lead to database
	createdLead, err := a.repo.CreateLead(ctx, lead)
	if err != nil {
		logger.Get().Error().Err(err).Msg("Failed to create lead")
		return nil, imhttp.NewCustomErr(http.StatusInternalServerError, "Failed to create lead", err.Error())
	}

	// Send OTP via SMS
	if err := a.smsClient.SendOTP(req.Phone, otp); err != nil {
		logger.Get().Error().Err(err).Msg("Failed to send OTP SMS")
		// Don't fail the lead creation if SMS fails, just log the error
	}

	// Send lead to external CRM
	go a.sendToCRM(ctx, createdLead, req)

	return &response.CreateLeadResponse{
		Message: "Leads Saved Successfully",
	}, nil
}

func (a *application) CreateLead(ctx context.Context, req *request.CreateLeadRequest) (*response.CreateLeadResponse, *imhttp.CustomError) {

	// Check for existing lead by phone number
	existingLead, err := a.repo.GetLeadByPhone(ctx, req.Phone)

	lead := &ent.Leads{
		Email:       req.Email,
		Name:        req.Name,
		Phone:       req.Phone,
		Message:     req.Message,
		Source:      "Organic",
		IsDuplicate: false,
		OtpVerified: false,
	}

	// Handle duplicate detection
	if err == nil && existingLead != nil {
		logger.Get().Info().Str("phone", req.Phone).Msg("Duplicate lead detected")
		lead.IsDuplicate = true
		lead.DuplicateReferenceID = strconv.Itoa(existingLead.ID)
	}

	// Set property or project relationship
	if req.PropertyID != "" {
		property, err := a.repo.GetPropertyByID(req.PropertyID)
		if err != nil {
			logger.Get().Warn().Err(err).Str("property_id", req.PropertyID).Msg("Property not found, creating lead without property association")
		} else {
			lead.Edges.Property = property
		}
	}

	if req.ProjectID != "" {
		project, err := a.repo.GetProjectByID(req.ProjectID)
		if err != nil {
			logger.Get().Warn().Err(err).Str("project_id", req.ProjectID).Msg("Project not found, creating lead without project association")
		} else {
			lead.Edges.Project = project
		}
	}

	// Save lead to database
	createdLead, err := a.repo.CreateLead(ctx, lead)
	if err != nil {
		logger.Get().Error().Err(err).Msg("Failed to create lead")
		return nil, imhttp.NewCustomErr(http.StatusInternalServerError, "Failed to create lead", err.Error())
	}

	// Send lead to external CRM
	go a.sendToCRM(ctx, createdLead, req)

	return &response.CreateLeadResponse{
		Message: "Leads Saved Successfully",
	}, nil
}

func (a *application) GetLeadByID(ctx context.Context, id int) (*response.Lead, *imhttp.CustomError) {
	lead, err := a.repo.GetLeadByID(ctx, id)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, imhttp.NewCustomErr(http.StatusNotFound, "Lead not found", "Lead not found")
		}
		return nil, imhttp.NewCustomErr(http.StatusInternalServerError, "Failed to get lead", err.Error())
	}

	return response.ToLeadResponse(lead), nil
}

func (a *application) GetAllLeads(ctx context.Context, req *request.GetLeadsRequest) (*response.DateLeadsData, *imhttp.CustomError) {
	return a.getLeadsByDateGrouped(ctx, req)
}

func (a *application) getLeadsByDateGrouped(ctx context.Context, req *request.GetLeadsRequest) (*response.DateLeadsData, *imhttp.CustomError) {
	// Build filters for all requests - support both single date and date range
	filters := make(map[string]interface{})

	if req.ProjectID != "" {
		filters["project_id"] = req.ProjectID
	}
	if req.PropertyID != "" {
		filters["property_id"] = req.PropertyID
	}
	if req.Phone != "" {
		filters["phone"] = req.Phone
	}
	if req.Source != "" {
		filters["source"] = req.Source
	}

	// Handle date filtering - support both single date and date range
	if req.Date != "" {
		filters["date"] = req.Date
	}
	if req.StartDate != "" {
		filters["start_date"] = req.StartDate
	}
	if req.EndDate != "" {
		filters["end_date"] = req.EndDate
	}

	// Get leads with all filters applied
	leads, err := a.repo.GetAllLeads(ctx, filters)
	if err != nil {
		return nil, imhttp.NewCustomErr(http.StatusInternalServerError, "Failed to get leads", err.Error())
	}

	// Separate unique and duplicate leads using sophisticated logic
	var uniqueLeads []*response.Lead
	duplicateGroups := make(map[string][]*response.Lead)

	for _, lead := range leads {
		leadResponse := response.ToLeadResponse(lead)

		if !lead.IsDuplicate {
			uniqueLeads = append(uniqueLeads, leadResponse)
		} else {
			// Group by duplicate_reference_id
			refID := lead.DuplicateReferenceID
			if refID == "" {
				// If no reference ID, treat as unique
				uniqueLeads = append(uniqueLeads, leadResponse)
			} else {
				duplicateGroups[refID] = append(duplicateGroups[refID], leadResponse)
			}
		}
	}

	// Process duplicate groups to create the required structure
	processedDuplicates := make(map[string]*response.DuplicateLeadGroup)
	for refID, duplicates := range duplicateGroups {
		if len(duplicates) == 0 {
			continue
		}

		// The leads are already ordered by created_at DESC from the repository
		// So the first one is the most recent (last)
		last := duplicates[0]
		var history []*response.Lead

		// Add remaining leads from the same date to history (already in reverse chronological order)
		if len(duplicates) > 1 {
			history = duplicates[1:]
		}

		// Build complete history chain including referenced leads from other dates
		referenceHistory, err := a.buildDuplicateHistory(ctx, refID)
		if err != nil {
			logger.Get().Error().Err(err).Str("reference_id", refID).Msg("Failed to build duplicate history")
		} else if len(referenceHistory) > 0 {
			// Append the reference history to the same-date history
			history = append(history, referenceHistory...)
		}

		processedDuplicates[refID] = &response.DuplicateLeadGroup{
			Last:    last,
			History: history,
		}
	}

	// Return the simplified response structure without date grouping
	return &response.DateLeadsData{
		UniqueLeads:    uniqueLeads,
		DuplicateLeads: processedDuplicates,
	}, nil
}

func (a *application) buildDuplicateHistory(ctx context.Context, referenceID string) ([]*response.Lead, error) {
	if referenceID == "" {
		return nil, nil
	}

	var history []*response.Lead
	currentRefID := referenceID

	for currentRefID != "" {
		refIDInt, err := strconv.Atoi(currentRefID)
		if err != nil {
			logger.Get().Warn().Str("reference_id", currentRefID).Msg("Invalid reference ID format")
			break
		}

		referencedLead, err := a.repo.GetLeadByID(ctx, refIDInt)
		if err != nil {
			if ent.IsNotFound(err) {
				logger.Get().Warn().Str("reference_id", currentRefID).Msg("Referenced lead not found")
			} else {
				logger.Get().Error().Err(err).Str("reference_id", currentRefID).Msg("Failed to get referenced lead")
			}
			break
		}

		history = append(history, response.ToLeadResponse(referencedLead))

		currentRefID = referencedLead.DuplicateReferenceID
	}

	return history, nil
}

func (a *application) ValidateOTP(ctx context.Context, req *request.ValidateOTPRequest) (*response.ValidateOTPResponse, *imhttp.CustomError) {
	lead, err := a.repo.GetLeadByPhoneAndOTP(ctx, req.Phone, req.OTP)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, imhttp.NewCustomErr(http.StatusBadRequest, "Invalid OTP or phone number", "Invalid OTP")
		}
		return nil, imhttp.NewCustomErr(http.StatusInternalServerError, "Failed to validate OTP", err.Error())
	}

	// Clear OTP and mark as verified
	lead.Otp = ""
	lead.OtpVerified = true

	_, err = a.repo.UpdateLead(ctx, lead)
	if err != nil {
		logger.Get().Error().Err(err).Msg("Failed to update lead after OTP validation")
		return nil, imhttp.NewCustomErr(http.StatusInternalServerError, "Failed to update lead", err.Error())
	}

	return &response.ValidateOTPResponse{
		Message: "OTP Validated Successfully",
	}, nil
}

func (a *application) ResendOTP(ctx context.Context, req *request.ResendOTPRequest) (*response.ResendOTPResponse, *imhttp.CustomError) {
	lead, err := a.repo.GetLeadByPhone(ctx, req.Phone)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, imhttp.NewCustomErr(http.StatusNotFound, "Lead not found", "Lead not found")
		}
		return nil, imhttp.NewCustomErr(http.StatusInternalServerError, "Failed to get lead", err.Error())
	}

	// Generate new OTP
	newOTP := generateOTP()
	lead.Otp = newOTP

	// Update lead with new OTP
	_, err = a.repo.UpdateLead(ctx, lead)
	if err != nil {
		logger.Get().Error().Err(err).Msg("Failed to update lead with new OTP")
		return nil, imhttp.NewCustomErr(http.StatusInternalServerError, "Failed to update lead", err.Error())
	}

	// Send new OTP via SMS
	if err := a.smsClient.SendOTP(req.Phone, newOTP); err != nil {
		logger.Get().Error().Err(err).Msg("Failed to send OTP SMS")
		return nil, imhttp.NewCustomErr(http.StatusInternalServerError, "Failed to send OTP", err.Error())
	}

	return &response.ResendOTPResponse{
		Message: "OTP Send Successfully",
	}, nil
}

// sendToCRM sends lead data to external CRM system asynchronously
func (a *application) sendToCRM(ctx context.Context, lead *ent.Leads, req *request.CreateLeadRequest) {
	projectName := ""

	// Get project name for CRM
	if lead.Edges.Project != nil {
		projectName = lead.Edges.Project.Name
	} else if lead.Edges.Property != nil && lead.Edges.Property.Edges.Project != nil {
		projectName = lead.Edges.Property.Edges.Project.Name
	}

	if projectName == "" {
		projectName = "Main Page"
	}

	crmData := client.CRMLeadData{
		Name:        lead.Name,
		Email:       lead.Email,
		Phone:       lead.Phone,
		ProjectName: projectName,
		QueryInfo:   lead.Message,
		Source:      lead.Source,
	}

	var syncStatus string
	if err := a.crmClient.SendLead(crmData); err != nil {
		logger.Get().Error().Err(err).Msg("Failed to send lead to CRM")
		syncStatus = "rejected"
	} else {
		logger.Get().Info().Int("lead_id", lead.ID).Msg("Lead sent to CRM successfully")
		syncStatus = "synced"
	}

	// Update sync status in database
	lead.SyncStatus = leads.SyncStatus(syncStatus)
	if _, err := a.repo.UpdateLead(context.Background(), lead); err != nil {
		logger.Get().Error().Err(err).Int("lead_id", lead.ID).Str("sync_status", syncStatus).Msg("Failed to update lead sync status")
	}
}
