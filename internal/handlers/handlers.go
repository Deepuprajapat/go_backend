package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
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
		fmt.Println("CLicked")
		err := json.NewDecoder(r.Body).Decode(&request)
		if err != nil {
			fmt.Println("ffrr")
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		_, err = strconv.ParseFloat(request.Phone, 64) // Can also use strconv.Atoi for integers
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		if request.OTP == ""{
			http.Error(w, err.Error(), http.StatusBadRequest)
			return 
		}
		// TODO: Implement token generation logic
		token := "generated_token"
		json.NewEncoder(w).Encode(map[string]string{"token": token})
	}
