package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

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

	return &imhttp.Response{
		StatusCode: http.StatusOK,
		Data:       result,
	}, nil
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
