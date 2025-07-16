package testutil

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/VI-IM/im_backend_go/ent"
	"github.com/VI-IM/im_backend_go/internal/application"
	"github.com/VI-IM/im_backend_go/internal/client"
	"github.com/VI-IM/im_backend_go/internal/config"
	"github.com/VI-IM/im_backend_go/internal/repository"
	"github.com/VI-IM/im_backend_go/internal/router"
	"github.com/VI-IM/im_backend_go/shared/logger"
	"github.com/gorilla/mux"
)

// TestServer represents a test server instance
type TestServer struct {
	Server      *httptest.Server
	Client      *ent.Client
	Config      config.Config
	Application application.ApplicationInterface
	cleanup     func()
}

// NewTestServer creates a new test server with database
func NewTestServer(t *testing.T) *TestServer {
	// Initialize logger for tests
	logger.Init()

	// Setup test database and client using testcontainers
	client, dbCleanup := SetupTestClient(t)

	// Setup test config (this will use the DATABASE_URL set by testcontainers)
	configCleanup := SetupTestConfig(t)
	cfg := config.GetConfig()

	// Create repositories
	repos := repository.NewRepository(client)

	// Create mock clients for testing
	mockS3Client := &mockS3Client{}
	mockSMSClient := &mockSMSClient{}
	mockCRMClient := &mockCRMClient{}

	// Create application layer
	app := application.NewApplication(repos, mockS3Client, mockSMSClient, mockCRMClient)

	// Initialize the global router with our application
	router.Init(app)

	// Create test server using the global router
	server := httptest.NewServer(router.Router)

	return &TestServer{
		Server:      server,
		Client:      client,
		Config:      cfg,
		Application: app,
		cleanup: func() {
			server.Close()
			configCleanup()
			dbCleanup()
		},
	}
}


// Close shuts down the test server and cleans up resources
func (ts *TestServer) Close() {
	ts.Server.Close()
	ts.cleanup()
}

// URL returns the base URL of the test server
func (ts *TestServer) URL() string {
	return ts.Server.URL
}

// NewRequest creates a new HTTP request against the test server
func (ts *TestServer) NewRequest(method, path string, body interface{}) (*http.Request, error) {
	url := fmt.Sprintf("%s%s", ts.URL(), path)
	return NewHTTPRequest(method, url, body)
}

// ClearDatabase truncates all tables for test isolation
func (ts *TestServer) ClearDatabase(t *testing.T) {
	ctx := context.Background()
	
	// Delete all data in order to respect foreign key constraints
	if _, err := ts.Client.User.Delete().Exec(ctx); err != nil {
		t.Fatalf("Failed to clear users: %v", err)
	}
	if _, err := ts.Client.Leads.Delete().Exec(ctx); err != nil {
		t.Fatalf("Failed to clear leads: %v", err)
	}
	if _, err := ts.Client.Property.Delete().Exec(ctx); err != nil {
		t.Fatalf("Failed to clear properties: %v", err)
	}
	if _, err := ts.Client.Project.Delete().Exec(ctx); err != nil {
		t.Fatalf("Failed to clear projects: %v", err)
	}
	if _, err := ts.Client.Developer.Delete().Exec(ctx); err != nil {
		t.Fatalf("Failed to clear developers: %v", err)
	}
	if _, err := ts.Client.Location.Delete().Exec(ctx); err != nil {
		t.Fatalf("Failed to clear locations: %v", err)
	}
	if _, err := ts.Client.Blogs.Delete().Exec(ctx); err != nil {
		t.Fatalf("Failed to clear blogs: %v", err)
	}
}

// Mock clients for testing

type mockS3Client struct{}

func (m *mockS3Client) UploadFile(ctx context.Context, key string, body io.Reader) (string, error) {
	return "https://test-bucket.s3.amazonaws.com/" + key, nil
}

func (m *mockS3Client) GetFile(ctx context.Context, key string) (io.ReadCloser, error) {
	return nil, fmt.Errorf("not implemented in mock")
}

func (m *mockS3Client) DeleteFile(ctx context.Context, key string) error {
	return nil
}

func (m *mockS3Client) GeneratePresignedURL(ctx context.Context, key string, operation string, duration time.Duration) (string, error) {
	return "https://test-presigned-url.com/" + key, nil
}

type mockSMSClient struct{}

func (m *mockSMSClient) SendOTP(phone, otp string) error {
	return nil
}

type mockCRMClient struct{}

func (m *mockCRMClient) SendLead(leadData client.CRMLeadData) error {
	return nil
}