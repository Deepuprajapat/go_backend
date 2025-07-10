package request

type GenerateTokenRequest struct {
	Email        string `json:"email"`
	Password     string `json:"password"`
	RefreshToken string `json:"refresh_token,omitempty"`
}

type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token"`
}

type SignupRequest struct {
	Username    string `json:"username" validate:"required"`
	Password    string `json:"password" validate:"required,min=6"`
	Email       string `json:"email" validate:"required,email"`
	Name        string `json:"name" validate:"required"`
	PhoneNumber string `json:"phone_number" validate:"required"`
}
