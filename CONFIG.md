# Configuration Guide

This application uses `envconfig` for loading configuration from environment variables with sensible defaults.

## Environment Variables

### Server Configuration

| Variable | Type | Default | Description |
|----------|------|---------|-------------|
| `PORT` | int | - | Server port |
| `HOST` | string | - | Server host |

### Database Configuration

| Variable | Type | Default | Description |
|----------|------|---------|-------------|
| `DB_PORT` | int | - | Database port |
| `DB_HOST` | string | - | Database host |
| `CONN_STRG` | string | - | Database connection string |

### JWT Configuration

| Variable | Type | Default | Description |
|----------|------|---------|-------------|
| `INVITATION_JWT_SECRET` | string | - | Secret for invitation JWT tokens |
| `AUTH_JWT_SECRET` | string | - | Secret for authentication JWT tokens |
| `JWT_EXPIRATION_DURATION` | string | - | JWT token expiration duration (e.g., "24h", "7d") |
| `JWT_COOKIE_NAME` | string | `"auth_token"` | Name of the authentication cookie |
| `JWT_COOKIE_DOMAIN` | string | `"localhost"` | Domain for the authentication cookie |
| `JWT_COOKIE_SECURE` | bool | `true` | Whether the cookie should only be sent over HTTPS |
| `JWT_COOKIE_HTTP_ONLY` | bool | `true` | Whether the cookie should be HTTP-only (not accessible via JavaScript) |

### Logging Configuration

| Variable | Type | Default | Description |
|----------|------|---------|-------------|
| `LOG_LEVEL` | string | `"2"` | Log level (debug, info, warn, error, fatal) |

### S3 Configuration

| Variable | Type | Default | Description |
|----------|------|---------|-------------|
| `S3_BUCKET` | string | - | S3 bucket name |
| `S3_REGION` | string | - | S3 region |
| `S3_ACCESS_KEY_ID` | string | - | S3 access key ID |
| `S3_SECRET_KEY` | string | - | S3 secret key |

## Usage

The configuration is automatically loaded when the application starts using:

```go
if err := config.LoadConfig(); err != nil {
    log.Fatal().Err(err).Msg("Failed to load configuration")
}

// Access configuration
cfg := config.GetConfig()
```

### Example Environment Setup

```bash
# Server
PORT=8080
HOST=0.0.0.0

# Database
DB_HOST=localhost
DB_PORT=5432
CONN_STRG=postgres://user:pass@localhost:5432/myapp?sslmode=disable

# JWT
AUTH_JWT_SECRET=your-super-secret-jwt-key-here
INVITATION_JWT_SECRET=your-invitation-secret-here
JWT_EXPIRATION_DURATION=168h
JWT_COOKIE_NAME=auth_token
JWT_COOKIE_DOMAIN=yourdomain.com
JWT_COOKIE_SECURE=true
JWT_COOKIE_HTTP_ONLY=true

# Logging
LOG_LEVEL=info

# S3 (optional)
S3_BUCKET=my-app-bucket
S3_REGION=us-east-1
S3_ACCESS_KEY_ID=your-access-key
S3_SECRET_KEY=your-secret-key
```

## Security Notes

1. **Always set JWT secrets** in production to strong, random values
2. Set `JWT_COOKIE_SECURE=true` when using HTTPS
3. Keep `JWT_COOKIE_HTTP_ONLY=true` to prevent XSS attacks
4. Never commit secrets to version control

## Error Handling

- Invalid boolean values will log a warning and use defaults
- Invalid integer values will log a warning and use defaults  
- Missing `CONN_STRG` when using `MustLoadDBConfig()` will cause the application to exit
- Missing `AUTH_JWT_SECRET` will show a warning when using `MustLoadConfig()` 