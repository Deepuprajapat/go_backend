package request

type GenerateTokenRequest struct {
	Username     string `json:"username"`
	Password     string `json:"password"`
	RefreshToken string `json:"refresh_token,omitempty"`
}

type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token"`
}
