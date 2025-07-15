package request

type CreateLeadRequest struct {
	PropertyID string `json:"property_id"`
	ProjectID  string `json:"project_id"`
	Name       string `json:"name" validate:"required"`
	Phone      string `json:"phone" validate:"required,len=10"`
	Email      string `json:"email" validate:"omitempty,email"`
	Message    string `json:"message"`
}

type ValidateOTPRequest struct {
	Phone string `json:"phone" validate:"required,len=10"`
	OTP   string `json:"otp" validate:"required,len=6"`
}

type ResendOTPRequest struct {
	Phone string `json:"phone" validate:"required,len=10"`
}

type GetLeadsRequest struct {
	ProjectID string `json:"project_id"`
	PropertyID string `json:"property_id"`
	Phone     string `json:"phone"`
	StartDate string `json:"start_date"`
	EndDate   string `json:"end_date"`
	Page      int    `json:"page"`
	Size      int    `json:"size"`
}