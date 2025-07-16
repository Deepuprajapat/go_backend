package testutil

import (
	"context"
	"os"
	"testing"

	"github.com/VI-IM/im_backend_go/ent"
	"github.com/VI-IM/im_backend_go/internal/auth"
	"github.com/VI-IM/im_backend_go/internal/utils"
	"github.com/google/uuid"
)

// TestUser represents a test user with credentials
type TestUser struct {
	ID       string
	Username string
	Email    string
	Password string
	Role     string
	Token    string
	Entity   *ent.User
}

// CreateTestUser creates a test user in the database
func CreateTestUser(t *testing.T, client *ent.Client, role string) *TestUser {
	// Set test JWT secret if not already set
	if os.Getenv("AUTH_JWT_SECRET") == "" {
		os.Setenv("AUTH_JWT_SECRET", "test-secret-key-for-testing-only")
	}

	ctx := context.Background()
	userID := uuid.New().String()
	username := "testuser_" + userID[:8]
	email := username + "@test.com"
	password := "testpassword123"

	// Hash password
	hashedPassword, err := utils.HashPassword(password)
	if err != nil {
		t.Fatalf("Failed to hash password: %v", err)
	}

	// Create user in database
	user, err := client.User.Create().
		SetID(userID).
		SetUsername(username).
		SetEmail(email).
		SetPassword(hashedPassword).
		SetName("Test User").
		SetRole(user.Role(role)).
		SetIsActive(true).
		SetIsVerified(true).
		Save(ctx)
	if err != nil {
		t.Fatalf("Failed to create test user: %v", err)
	}

	// Generate JWT token
	isAdmin := role == "superadmin"
	token, err := auth.GenerateToken(userID, isAdmin, role, "")
	if err != nil {
		t.Fatalf("Failed to generate token: %v", err)
	}

	return &TestUser{
		ID:       userID,
		Username: username,
		Email:    email,
		Password: password,
		Role:     role,
		Token:    token,
		Entity:   user,
	}
}

// CreateBusinessPartner creates a test business partner user
func CreateBusinessPartner(t *testing.T, client *ent.Client) *TestUser {
	return CreateTestUser(t, client, "business_partner")
}

// CreateSuperAdmin creates a test super admin user
func CreateSuperAdmin(t *testing.T, client *ent.Client) *TestUser {
	return CreateTestUser(t, client, "superadmin")
}

// CreateDM creates a test DM user
func CreateDM(t *testing.T, client *ent.Client) *TestUser {
	return CreateTestUser(t, client, "dm")
}

// AuthorizedClient creates an HTTP client with authentication token
func (tu *TestUser) AuthorizedClient(baseURL string) *HTTPClient {
	client := NewHTTPClient(baseURL)
	client.SetAuthToken(tu.Token)
	return client
}

// LoginRequest represents a login request body
type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// LoginResponse represents a login response
type LoginResponse struct {
	Token     string      `json:"token"`
	User      interface{} `json:"user"`
	Message   string      `json:"message"`
	IsSuccess bool        `json:"is_success"`
}

// Login performs a login request and returns the response
func (tu *TestUser) Login(client *HTTPClient) (*LoginResponse, error) {
	loginReq := LoginRequest{
		Username: tu.Username,
		Password: tu.Password,
	}

	var loginResp LoginResponse
	_, err := client.POST("/v1/api/auth/login", loginReq, &loginResp)
	if err != nil {
		return nil, err
	}

	// Update token from response
	if loginResp.Token != "" {
		tu.Token = loginResp.Token
		client.SetAuthToken(tu.Token)
	}

	return &loginResp, nil
}