package config

import (
	"log"
	"strings"
	"time"

	"github.com/kelseyhightower/envconfig"
	"github.com/rs/zerolog"
)

var cfg Config

type (
	Config struct {
		Server
		Database
		JWTConfig
		Logger
		S3
		CRM
		TLS
	}

	Server struct {
		Port    int    `envconfig:"PORT"`
		Host    string `envconfig:"HOST"`
		BaseURL string `envconfig:"BASE_URL" default:"https://investmango.com"`
	}
	
	TLS struct {
		Enabled      bool   `envconfig:"TLS_ENABLED" default:"false"`
		CertPath     string `envconfig:"TLS_CERT_PATH"`
		KeyPath      string `envconfig:"TLS_KEY_PATH"`
		Port         int    `envconfig:"TLS_PORT" default:"8443"`
		HTTPSOnly    bool   `envconfig:"HTTPS_ONLY" default:"false"`
		AutoCert     bool   `envconfig:"TLS_AUTO_CERT" default:"false"`
		AutoCertHost string `envconfig:"TLS_AUTO_CERT_HOST"`
	}
	Database struct {
		DB_Port int    `envconfig:"DB_PORT"`
		DB_HOST string `envconfig:"DB_HOST"`
		URL     string `envconfig:"DATABASE_URL"`
	}
	JWTConfig struct {
		AuthSecret        string        `envconfig:"AUTH_JWT_SECRET"`
		ExpiresIn         string        `envconfig:"JWT_EXPIRATION_DURATION" default:"24h"`
		ExpiresInDuration time.Duration `envconfig:"-"`
	}
	Logger struct {
		Level zerolog.Level `envconfig:"LOG_LEVEL" default:"1"`
		Mode  string        `envconfig:"LOG_MODE" default:"simple"`
	}

	S3 struct {
		Bucket      string `envconfig:"AWS_BUCKET"`
		Region      string `envconfig:"AWS_REGION"`
		AccessKeyID string `envconfig:"AWS_ACCESS_KEY_ID"`
		SecretKey   string `envconfig:"AWS_SECRET_KEY"`
	}

	CRM struct {
		BaseURL    string        `envconfig:"CRM_BASE_URL" default:"http://148.66.133.154:8181/new-leads/from/open-source"`
		Timeout    time.Duration `envconfig:"CRM_TIMEOUT" default:"30s"`
		Enabled    bool          `envconfig:"CRM_ENABLED" default:"true"`
		MaxRetries int           `envconfig:"CRM_MAX_RETRIES" default:"3"`
	}
)

func LoadConfig() error {
	if err := envconfig.Process("", &cfg); err != nil {
		log.Fatalf("error processing config: %v", err)
		return err
	}

	drtn, err := time.ParseDuration(cfg.ExpiresIn)
	if err != nil {
		log.Fatalf("error parsing duration: %v", err)
		return err
	}

	cfg.ExpiresInDuration = drtn

	return nil
}

func GetConfig() Config {
	return cfg
}

func (l *Logger) GetLogLevel() zerolog.Level {
	switch strings.ToLower(l.Level.String()) {
	case "debug":
		return zerolog.DebugLevel
	case "info":
		return zerolog.InfoLevel
	case "warn":
		return zerolog.WarnLevel
	case "error":
		return zerolog.ErrorLevel
	case "fatal":
		return zerolog.FatalLevel
	default:
		return zerolog.InfoLevel
	}
}
