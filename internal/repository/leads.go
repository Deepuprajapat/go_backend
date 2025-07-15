package repository

import (
	"context"
	"time"

	"github.com/VI-IM/im_backend_go/ent"
	"github.com/VI-IM/im_backend_go/ent/leads"
	"github.com/VI-IM/im_backend_go/ent/project"
	"github.com/VI-IM/im_backend_go/ent/property"
	"github.com/VI-IM/im_backend_go/shared/logger"
)

func (r *repository) CreateLead(ctx context.Context, lead *ent.Leads) (*ent.Leads, error) {
	logger.Get().Info().
		Str("phone", lead.Phone).
		Str("email", lead.Email).
		Msg("Creating new lead")

	leadBuilder := r.db.Leads.Create().
		SetName(lead.Name).
		SetPhone(lead.Phone).
		SetSource(lead.Source).
		SetIsDuplicate(lead.IsDuplicate).
		SetOtpVerified(lead.OtpVerified)
	
	if lead.Email != "" {
		leadBuilder.SetEmail(lead.Email)
	}

	if lead.Otp != "" {
		leadBuilder.SetOtp(lead.Otp)
	}

	if lead.Message != "" {
		leadBuilder.SetMessage(lead.Message)
	}

	if lead.DuplicateReferenceID != "" {
		leadBuilder.SetDuplicateReferenceID(lead.DuplicateReferenceID)
	}

	if lead.Edges.Property != nil {
		leadBuilder.SetProperty(lead.Edges.Property)
	}

	if lead.Edges.Project != nil {
		leadBuilder.SetProject(lead.Edges.Project)
	}

	createdLead, err := leadBuilder.Save(ctx)
	if err != nil {
		logger.Get().Error().Err(err).Msg("Failed to create lead")
		return nil, err
	}

	logger.Get().Info().
		Int("lead_id", createdLead.ID).
		Msg("Lead created successfully")

	return createdLead, nil
}

func (r *repository) GetLeadByID(ctx context.Context, id int) (*ent.Leads, error) {
	lead, err := r.db.Leads.Query().
		Where(leads.ID(id)).
		WithProperty().
		WithProject().
		First(ctx)

	if err != nil {
		if ent.IsNotFound(err) {
			logger.Get().Debug().Int("id", id).Msg("Lead not found")
			return nil, err
		}
		logger.Get().Error().Err(err).Int("id", id).Msg("Failed to get lead by ID")
		return nil, err
	}

	return lead, nil
}

func (r *repository) GetLeadByPhone(ctx context.Context, phone string) (*ent.Leads, error) {
	lead, err := r.db.Leads.Query().
		Where(leads.Phone(phone)).
		WithProperty().
		WithProject().
		Order(ent.Desc(leads.FieldCreatedAt)).
		First(ctx)

	if err != nil {
		if ent.IsNotFound(err) {
			logger.Get().Debug().Str("phone", phone).Msg("Lead not found")
			return nil, err
		}
		logger.Get().Error().Err(err).Str("phone", phone).Msg("Failed to get lead by phone")
		return nil, err
	}

	return lead, nil
}

func (r *repository) GetLeadByPhoneAndOTP(ctx context.Context, phone, otp string) (*ent.Leads, error) {
	lead, err := r.db.Leads.Query().
		Where(
			leads.Phone(phone),
			leads.Otp(otp),
		).
		WithProperty().
		WithProject().
		First(ctx)

	if err != nil {
		if ent.IsNotFound(err) {
			logger.Get().Debug().Str("phone", phone).Msg("Lead not found with OTP")
			return nil, err
		}
		logger.Get().Error().Err(err).Str("phone", phone).Msg("Failed to get lead by phone and OTP")
		return nil, err
	}

	return lead, nil
}

func (r *repository) UpdateLead(ctx context.Context, lead *ent.Leads) (*ent.Leads, error) {
	updateBuilder := r.db.Leads.UpdateOneID(lead.ID).
		SetUpdatedAt(time.Now())

	if lead.Otp != "" {
		updateBuilder.SetOtp(lead.Otp)
	} else {
		updateBuilder.ClearOtp()
	}

	if lead.OtpVerified {
		updateBuilder.SetOtpVerified(lead.OtpVerified)
	}

	if lead.IsDuplicate {
		updateBuilder.SetIsDuplicate(lead.IsDuplicate)
	}

	if lead.DuplicateReferenceID != "" {
		updateBuilder.SetDuplicateReferenceID(lead.DuplicateReferenceID)
	}

	updatedLead, err := updateBuilder.Save(ctx)
	if err != nil {
		logger.Get().Error().Err(err).Int("lead_id", lead.ID).Msg("Failed to update lead")
		return nil, err
	}

	logger.Get().Info().Int("lead_id", lead.ID).Msg("Lead updated successfully")
	return updatedLead, nil
}

func (r *repository) GetAllLeads(ctx context.Context, offset, limit int, filters map[string]interface{}) ([]*ent.Leads, int, error) {
	query := r.db.Leads.Query().
		WithProperty().
		WithProject()

	// Apply filters
	if projectID, ok := filters["project_id"].(string); ok && projectID != "" {
		query = query.Where(leads.HasProjectWith(project.ID(projectID)))
	}

	if propertyID, ok := filters["property_id"].(string); ok && propertyID != "" {
		query = query.Where(leads.HasPropertyWith(property.ID(propertyID)))
	}

	if phone, ok := filters["phone"].(string); ok && phone != "" {
		query = query.Where(leads.Phone(phone))
	}

	if startDate, ok := filters["start_date"].(string); ok && startDate != "" {
		if startTime, err := time.Parse(time.RFC3339, startDate); err == nil {
			query = query.Where(leads.CreatedAtGTE(startTime))
		}
	}

	if endDate, ok := filters["end_date"].(string); ok && endDate != "" {
		if endTime, err := time.Parse(time.RFC3339, endDate); err == nil {
			query = query.Where(leads.CreatedAtLTE(endTime))
		}
	}

	// Get total count
	totalCount, err := query.Count(ctx)
	if err != nil {
		logger.Get().Error().Err(err).Msg("Failed to count leads")
		return nil, 0, err
	}

	// Apply pagination and ordering
	leadsData, err := query.
		Order(ent.Desc(leads.FieldCreatedAt)).
		Offset(offset).
		Limit(limit).
		All(ctx)

	if err != nil {
		logger.Get().Error().Err(err).Msg("Failed to get leads")
		return nil, 0, err
	}

	logger.Get().Info().
		Int("total_count", totalCount).
		Int("returned_count", len(leadsData)).
		Msg("Retrieved leads successfully")

	return leadsData, totalCount, nil
}