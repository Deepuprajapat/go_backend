package testutil

import (
	"os"
	"testing"

	"github.com/VI-IM/im_backend_go/internal/config"
)

// SetupTestConfig sets up test environment variables and configuration
func SetupTestConfig(t *testing.T) func() {
	// Store original environment variables
	originalVars := make(map[string]string)
	testVars := map[string]string{
		"PORT":                      "8080",
		"HOST":                      "localhost",
		// DATABASE_URL is set by testcontainers, don't override it
		"AUTH_JWT_SECRET":           "test-secret-key-for-testing-only-do-not-use-in-production",
		"JWT_EXPIRATION_DURATION":   "24h",
		"LOG_LEVEL":                 "debug",
		"LOG_MODE":                  "simple",
		"AWS_BUCKET":                "test-bucket",
		"AWS_REGION":                "us-east-1",
		"AWS_ACCESS_KEY_ID":         "test-access-key",
		"AWS_SECRET_KEY":            "test-secret-key",
	}

	// Also backup DATABASE_URL even though we don't set it (testcontainers will)
	originalVars["DATABASE_URL"] = os.Getenv("DATABASE_URL")

	// Backup and set test environment variables
	for key, value := range testVars {
		originalVars[key] = os.Getenv(key)
		os.Setenv(key, value)
	}

	// Load config with test values
	if err := config.LoadConfig(); err != nil {
		t.Fatalf("Failed to load test config: %v", err)
	}

	// Return cleanup function
	return func() {
		// Restore original environment variables
		for key, originalValue := range originalVars {
			if originalValue == "" {
				os.Unsetenv(key)
			} else {
				os.Setenv(key, originalValue)
			}
		}
	}
}

// TestConfig returns a test configuration
func TestConfig() config.Config {
	return config.Config{
		Server: config.Server{
			Port: 8080,
			Host: "localhost",
		},
		Database: config.Database{
			URL: "postgres://postgres:password@localhost:5432/im_test_db?sslmode=disable",
		},
		JWTConfig: config.JWTConfig{
			AuthSecret: "test-secret-key-for-testing-only-do-not-use-in-production",
			ExpiresIn:  "24h",
		},
		S3: config.S3{
			Bucket:      "test-bucket",
			Region:      "us-east-1",
			AccessKeyID: "test-access-key",
			SecretKey:   "test-secret-key",
		},
	}
}

// IsTestEnvironment checks if we're running in test environment
func IsTestEnvironment() bool {
	return os.Getenv("GO_ENV") == "test" || 
		   os.Getenv("AUTH_JWT_SECRET") == "test-secret-key-for-testing-only-do-not-use-in-production"
}

// SetTestEnv marks the environment as test
func SetTestEnv() {
	os.Setenv("GO_ENV", "test")
}

// UnsetTestEnv removes test environment marker
func UnsetTestEnv() {
	os.Unsetenv("GO_ENV")
}