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

var istLocation *time.Location

func init() {
	var err error
	istLocation, err = time.LoadLocation("Asia/Kolkata")
	if err != nil {
		istLocation = time.FixedZone("IST", 5*60*60+30*60)
	}
}

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
		WithProperty(func(q *ent.PropertyQuery) {
			q.WithProject()
		}).
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
		WithProperty(func(q *ent.PropertyQuery) {
			q.WithProject()
		}).
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
		WithProperty(func(q *ent.PropertyQuery) {
			q.WithProject()
		}).
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

	if lead.SyncStatus != "" {
		updateBuilder.SetSyncStatus(leads.SyncStatus(lead.SyncStatus))
	}

	updatedLead, err := updateBuilder.Save(ctx)
	if err != nil {
		logger.Get().Error().Err(err).Int("lead_id", lead.ID).Msg("Failed to update lead")
		return nil, err
	}

	logger.Get().Info().Int("lead_id", lead.ID).Msg("Lead updated successfully")
	return updatedLead, nil
}

func (r *repository) GetAllLeads(ctx context.Context, filters map[string]interface{}) ([]*ent.Leads, error) {
	query := r.db.Leads.Query().
		WithProperty(func(q *ent.PropertyQuery) {
			q.WithProject()
		}).
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

	if source, ok := filters["source"].(string); ok && source != "" {
		query = query.Where(leads.Source(source))
	}

	if startDate, ok := filters["start_date"].(string); ok && startDate != "" {
		if startTime, err := time.Parse(time.DateOnly, startDate); err == nil {
			// Convert to IST and set to start of day (00:00:00)
			startTimeIST := time.Date(startTime.Year(), startTime.Month(), startTime.Day(), 0, 0, 0, 0, istLocation)
			query = query.Where(leads.CreatedAtGTE(startTimeIST))
		}
	}

	if endDate, ok := filters["end_date"].(string); ok && endDate != "" {
		if endTime, err := time.Parse(time.DateOnly, endDate); err == nil {
			// Convert to IST and set to end of day (23:59:59)
			endTimeIST := time.Date(endTime.Year(), endTime.Month(), endTime.Day(), 23, 59, 59, 999999999, istLocation)
			query = query.Where(leads.CreatedAtLTE(endTimeIST))
		}
	}

	// Handle specific date filter (format: 2025-07-16)
	if date, ok := filters["date"].(string); ok && date != "" {
		if dateTime, err := time.Parse(time.DateOnly, date); err == nil {
			// Get start and end of the day in IST
			startOfDay := time.Date(dateTime.Year(), dateTime.Month(), dateTime.Day(), 0, 0, 0, 0, istLocation)
			endOfDay := time.Date(dateTime.Year(), dateTime.Month(), dateTime.Day(), 23, 59, 59, 999999999, istLocation)
			query = query.Where(
				leads.CreatedAtGTE(startOfDay),
				leads.CreatedAtLTE(endOfDay),
			)
		}
	}

	// Apply ordering and get all results
	leadsData, err := query.
		Order(ent.Desc(leads.FieldCreatedAt)).
		All(ctx)

	if err != nil {
		logger.Get().Error().Err(err).Msg("Failed to get leads")
		return nil, err
	}

	logger.Get().Info().
		Int("count", len(leadsData)).
		Msg("Retrieved leads successfully")

	return leadsData, nil
}

func (r *repository) GetLeadsByDate(ctx context.Context, date string) ([]*ent.Leads, error) {
	dateTime, err := time.Parse("2006-01-02", date)
	if err != nil {
		logger.Get().Error().Err(err).Str("date", date).Msg("Invalid date format")
		return nil, err
	}

	// Get start and end of the day in IST
	startOfDay := time.Date(dateTime.Year(), dateTime.Month(), dateTime.Day(), 0, 0, 0, 0, istLocation)
	endOfDay := time.Date(dateTime.Year(), dateTime.Month(), dateTime.Day(), 23, 59, 59, 999999999, istLocation)

	leadsData, err := r.db.Leads.Query().
		Where(
			leads.CreatedAtGTE(startOfDay),
			leads.CreatedAtLTE(endOfDay),
		).
		WithProperty(func(q *ent.PropertyQuery) {
			q.WithProject()
		}).
		WithProject().
		Order(ent.Desc(leads.FieldCreatedAt)).
		All(ctx)

	if err != nil {
		logger.Get().Error().Err(err).Str("date", date).Msg("Failed to get leads by date")
		return nil, err
	}

	logger.Get().Info().
		Str("date", date).
		Int("count", len(leadsData)).
		Msg("Retrieved leads by date successfully")

	return leadsData, nil
}
