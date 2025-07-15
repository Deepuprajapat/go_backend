package application

import (
	"context"
	"fmt"
	"math/rand"
	"net/http"
	"strconv"
	"time"

	"github.com/VI-IM/im_backend_go/ent"
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

func (a *application) GetAllLeads(ctx context.Context, req *request.GetLeadsRequest) (*response.LeadListResponse, *imhttp.CustomError) {
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
	if req.StartDate != "" {
		filters["start_date"] = req.StartDate
	}
	if req.EndDate != "" {
		filters["end_date"] = req.EndDate
	}

	// Set default pagination values
	if req.Size <= 0 {
		req.Size = 12
	}
	if req.Page < 0 {
		req.Page = 0
	}

	offset := req.Page * req.Size

	leads, totalCount, err := a.repo.GetAllLeads(ctx, offset, req.Size, filters)
	if err != nil {
		return nil, imhttp.NewCustomErr(http.StatusInternalServerError, "Failed to get leads", err.Error())
	}

	return response.ToLeadListResponse(leads, totalCount, req.Size, req.Page), nil
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

	crmData := client.CRMLeadData{
		Name:        lead.Name,
		Email:       lead.Email,
		Phone:       lead.Phone,
		ProjectName: projectName,
		QueryInfo:   lead.Message,
		Source:      lead.Source,
	}

	if err := a.crmClient.SendLead(crmData); err != nil {
		logger.Get().Error().Err(err).Msg("Failed to send lead to CRM")
	}
}