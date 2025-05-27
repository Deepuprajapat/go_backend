package handlers

import (
	"encoding/json"
	"net/http"
)

// HealthCheck handles the health check endpoint
func HealthCheck(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

// GenerateToken handles token generation requests
func GenerateToken(w http.ResponseWriter, r *http.Request) {
	var request struct {
		Phone string `json:"phone"`
		OTP   string `json:"otp"`
	}
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// TODO: Implement token generation logic
	token := "generated_token"
	json.NewEncoder(w).Encode(map[string]string{"token": token})
}
