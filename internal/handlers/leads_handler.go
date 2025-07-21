package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/VI-IM/im_backend_go/internal/auth"
	"github.com/VI-IM/im_backend_go/request"
	imhttp "github.com/VI-IM/im_backend_go/shared"
	"github.com/VI-IM/im_backend_go/shared/logger"
	"github.com/gorilla/mux"
)

func (h *Handler) CreateLeadWithOTP(r *http.Request) (*imhttp.Response, *imhttp.CustomError) {
	var req request.CreateLeadRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		logger.Get().Error().Err(err).Msg("Failed to decode create lead request")
		return nil, imhttp.NewCustomErr(http.StatusBadRequest, "Invalid request body", err.Error())
	}

	// Validate request
	if req.Name == "" {
		return nil, imhttp.NewCustomErr(http.StatusBadRequest, "Name is required", "Validation error")
	}
	if req.Phone == "" || len(req.Phone) != 10 {
		return nil, imhttp.NewCustomErr(http.StatusBadRequest, "Phone must be 10 digits", "Validation error")
	}

	result, customErr := h.app.CreateLeadWithOTP(r.Context(), &req)
	if customErr != nil {
		return nil, customErr
	}

	return &imhttp.Response{
		StatusCode: http.StatusOK,
		Data:       result,
	}, nil
}

func (h *Handler) CreateLead(r *http.Request) (*imhttp.Response, *imhttp.CustomError) {
	var req request.CreateLeadRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		logger.Get().Error().Err(err).Msg("Failed to decode create lead request")
		return nil, imhttp.NewCustomErr(http.StatusBadRequest, "Invalid request body", err.Error())
	}

	// Validate request
	if req.Name == "" {
		return nil, imhttp.NewCustomErr(http.StatusBadRequest, "Name is required", "Validation error")
	}
	if req.Phone == "" || len(req.Phone) != 10 {
		return nil, imhttp.NewCustomErr(http.StatusBadRequest, "Phone must be 10 digits", "Validation error")
	}

	result, customErr := h.app.CreateLead(r.Context(), &req)
	if customErr != nil {
		return nil, customErr
	}

	return &imhttp.Response{
		StatusCode: http.StatusOK,
		Data:       result,
	}, nil
}

func (h *Handler) GetLeadByID(r *http.Request) (*imhttp.Response, *imhttp.CustomError) {
	vars := mux.Vars(r)
	idStr := vars["id"]

	id, err := strconv.Atoi(idStr)
	if err != nil {
		return nil, imhttp.NewCustomErr(http.StatusBadRequest, "Invalid ID format", "ID must be a number")
	}

	result, customErr := h.app.GetLeadByID(r.Context(), id)
	if customErr != nil {
		return nil, customErr
	}

	// Validate property ownership for business partners
	claims, ok := r.Context().Value("user_claims").(*auth.Claims)
	if ok && claims.Role == "business_partner" {
		// If lead has a property, validate ownership
		if result.PropertyID != "" {
			if err := h.validatePropertyOwnership(r, []string{result.PropertyID}); err != nil {
				return nil, err
			}
		}
	}

	return &imhttp.Response{
		StatusCode: http.StatusOK,
		Data:       result,
	}, nil
}

// validatePropertyOwnership validates if a business partner can access leads for the given property IDs
func (h *Handler) validatePropertyOwnership(r *http.Request, propertyIDs []string) *imhttp.CustomError {
	claims, ok := r.Context().Value("user_claims").(*auth.Claims)
	if !ok {
		return imhttp.NewCustomErr(http.StatusUnauthorized, "Invalid user context", "Invalid user context")
	}

	// Superadmin and dm can access all leads
	if claims.Role == "superadmin" || claims.Role == "dm" {
		return nil
	}

	// Business partner can only access leads for properties they created
	if claims.Role == "business_partner" {
		for _, propertyID := range propertyIDs {
			if propertyID == "" {
				continue
			}
			property, err := h.app.GetPropertyByID(propertyID)
			if err != nil {
				logger.Get().Error().Err(err).Str("property_id", propertyID).Msg("Failed to get property for ownership check")
				return imhttp.NewCustomErr(http.StatusInternalServerError, "Failed to verify property ownership", err.Error())
			}
			// Check if the property was created by this user
			if property.CreatedByUserID == "" || property.CreatedByUserID != claims.UserID {
				return imhttp.NewCustomErr(http.StatusForbidden, "Access denied: you can only access leads for properties you created", "Property ownership check failed")
			}
		}
	}

	return nil
}

func (h *Handler) GetAllLeads(r *http.Request) (*imhttp.Response, *imhttp.CustomError) {
	queryParams := r.URL.Query()

	var req request.GetLeadsRequest
	req.ProjectID = queryParams.Get("project_id")
	req.PropertyID = queryParams.Get("property_id")
	req.Phone = queryParams.Get("phone")
	req.StartDate = queryParams.Get("start_date")
	req.EndDate = queryParams.Get("end_date")
	req.Date = queryParams.Get("date")
	req.Source = queryParams.Get("source")

	// Handle multiple property IDs from query parameter
	if propertyIDsParam := queryParams.Get("property_ids"); propertyIDsParam != "" {
		req.PropertyIDs = strings.Split(propertyIDsParam, ",")
		// Trim whitespace from each property ID
		for i, id := range req.PropertyIDs {
			req.PropertyIDs[i] = strings.TrimSpace(id)
		}
	}

	// If single property_id is provided, add it to the PropertyIDs slice
	if req.PropertyID != "" {
		req.PropertyIDs = append(req.PropertyIDs, req.PropertyID)
	}

	// Validate property ownership for business partners
	if len(req.PropertyIDs) > 0 {
		if err := h.validatePropertyOwnership(r, req.PropertyIDs); err != nil {
			return nil, err
		}
	}

	result, customErr := h.app.GetAllLeads(r.Context(), &req)
	if customErr != nil {
		return nil, customErr
	}

	return &imhttp.Response{
		StatusCode: http.StatusOK,
		Data:       result,
	}, nil
}

func (h *Handler) ValidateOTP(r *http.Request) (*imhttp.Response, *imhttp.CustomError) {
	queryParams := r.URL.Query()

	phone := queryParams.Get("phone")
	otp := queryParams.Get("OTP")

	if phone == "" {
		return nil, imhttp.NewCustomErr(http.StatusBadRequest, "Phone parameter is required", "Validation error")
	}
	if otp == "" {
		return nil, imhttp.NewCustomErr(http.StatusBadRequest, "OTP parameter is required", "Validation error")
	}
	if len(phone) != 10 {
		return nil, imhttp.NewCustomErr(http.StatusBadRequest, "Phone must be 10 digits", "Validation error")
	}
	if len(otp) != 6 {
		return nil, imhttp.NewCustomErr(http.StatusBadRequest, "OTP must be 6 digits", "Validation error")
	}

	req := &request.ValidateOTPRequest{
		Phone: phone,
		OTP:   otp,
	}

	result, customErr := h.app.ValidateOTP(r.Context(), req)
	if customErr != nil {
		return nil, customErr
	}

	return &imhttp.Response{
		StatusCode: http.StatusOK,
		Data:       result,
	}, nil
}

func (h *Handler) ResendOTP(r *http.Request) (*imhttp.Response, *imhttp.CustomError) {
	queryParams := r.URL.Query()

	phone := queryParams.Get("phone")

	if phone == "" {
		return nil, imhttp.NewCustomErr(http.StatusBadRequest, "Phone parameter is required", "Validation error")
	}
	if len(phone) != 10 {
		return nil, imhttp.NewCustomErr(http.StatusBadRequest, "Phone must be 10 digits", "Validation error")
	}

	req := &request.ResendOTPRequest{
		Phone: phone,
	}

	result, customErr := h.app.ResendOTP(r.Context(), req)
	if customErr != nil {
		return nil, customErr
	}

	return &imhttp.Response{
		StatusCode: http.StatusOK,
		Data:       result,
	}, nil
}
