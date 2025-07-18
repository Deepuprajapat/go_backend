package request

type AddLocationRequest struct {
	LocalityName string `json:"locality_name" validate:"required"`
	City         string `json:"city" validate:"required"`
	State        string `json:"state" validate:"required"`
	PhoneNumber  string `json:"phone_number,omitempty"`
	Country      string `json:"country" validate:"required"`
	Pincode      string `json:"pincode,omitempty"`
}
