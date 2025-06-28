package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/VI-IM/im_backend_go/internal/application"
	"github.com/VI-IM/im_backend_go/request"
	imhttp "github.com/VI-IM/im_backend_go/shared"
	"github.com/VI-IM/im_backend_go/shared/logger"
	"github.com/gorilla/mux"
)

type PropertyHandler struct {
	app application.ApplicationInterface
}

func NewPropertyHandler(app application.ApplicationInterface) *PropertyHandler {
	return &PropertyHandler{app: app}
}

func (h *PropertyHandler) GetProperty(r *http.Request) (*imhttp.Response, *imhttp.CustomError) {
	vars := mux.Vars(r)
	propertyID := vars["property_id"]

	response, err := h.app.GetPropertyByID(propertyID)
	if err != nil {
		logger.Get().Error().Err(err).Msg("Failed to get property")
		return nil, imhttp.NewCustomErr(http.StatusInternalServerError, "Failed to get property", err.Error())
	}

	return &imhttp.Response{
		Data:       response,
		StatusCode: http.StatusOK,
	}, nil
}

func (h *PropertyHandler) UpdateProperty(r *http.Request) (*imhttp.Response, *imhttp.CustomError) {
	vars := mux.Vars(r)
	propertyID := vars["property_id"]
	if propertyID == "" {
		logger.Get().Error().Msg("Property ID is required")
		return nil, imhttp.NewCustomErr(http.StatusBadRequest, "Property ID is required", "Property ID is required")
	}

	var input request.UpdatePropertyRequest
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		logger.Get().Error().Msg("Invalid request body")
		return nil, imhttp.NewCustomErr(http.StatusBadRequest, "Invalid request body", err.Error())
	}

	input.PropertyID = propertyID
	response, err := h.app.UpdateProperty(input)
	if err != nil {
		return nil, imhttp.NewCustomErr(http.StatusInternalServerError, "Failed to update property", err.Error())
	}

	return &imhttp.Response{
		Data:       response,
		StatusCode: http.StatusOK,
	}, nil
}

func (h *PropertyHandler) GetPropertiesOfProject(r *http.Request) (*imhttp.Response, *imhttp.CustomError) {
	vars := mux.Vars(r)
	projectID := vars["project_id"]
	if projectID == "" {
		logger.Get().Error().Msg("Project ID is required")
		return nil, imhttp.NewCustomErr(http.StatusBadRequest, "Project ID is required", "Project ID is required")
	}

	response, err := h.app.GetPropertiesOfProject(projectID)
	if err != nil {
		return nil, imhttp.NewCustomErr(http.StatusInternalServerError, "Failed to get properties of project", err.Error())
	}

	return &imhttp.Response{
		Data:       response,
		StatusCode: http.StatusOK,
	}, nil
}

func (h *PropertyHandler) AddProperty(r *http.Request) (*imhttp.Response, *imhttp.CustomError) {
	var input request.AddPropertyRequest
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		logger.Get().Error().Msg("Invalid request body")
		return nil, imhttp.NewCustomErr(http.StatusBadRequest, "Invalid request body", err.Error())
	}

	if input.ProjectID == "" || input.PropertyType == "" || input.AgeOfProperty == "" || input.FloorNumber == "" || input.Facing == "" || input.Furnishing == "" || input.BalconyCount == "" {
		logger.Get().Error().Msg("Invalid request body")
		return nil, imhttp.NewCustomErr(http.StatusBadRequest, "Invalid request body", "Invalid request body")
	}

	propertyID, err := h.app.AddProperty(input)
	if err != nil {
		return nil, imhttp.NewCustomErr(http.StatusInternalServerError, "Failed to add property", err.Error())
	}

	return &imhttp.Response{
		Data:       propertyID,
		StatusCode: http.StatusOK,
	}, nil
}
