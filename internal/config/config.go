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
	}

	Server struct {
		Port int    `envconfig:"PORT"`
		Host string `envconfig:"HOST"`
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
