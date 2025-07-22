# TLS Certificate Management

This directory contains TLS certificates for HTTPS support. The certificates are excluded from version control for security reasons.

## Development Setup

### Option 1: Self-Signed Certificate (Development Only)

Generate a self-signed certificate for local development:

```bash
# Generate private key
openssl genrsa -out server.key 2048

# Generate certificate signing request
openssl req -new -key server.key -out server.csr

# Generate self-signed certificate (valid for 365 days)
openssl x509 -req -days 365 -in server.csr -signkey server.key -out server.crt

# Clean up CSR file
rm server.csr
```

### Option 2: mkcert (Recommended for Development)

Install mkcert for locally-trusted development certificates:

```bash
# macOS
brew install mkcert
mkcert -install

# Generate certificate for localhost
mkcert -key-file server.key -cert-file server.crt localhost 127.0.0.1 ::1
```

## Production Setup

### Option 1: Let's Encrypt (Automatic)

Enable automatic certificate management:

```bash
export TLS_ENABLED=true
export TLS_AUTO_CERT=true
export TLS_AUTO_CERT_HOST=yourdomain.com
```

The server will automatically obtain and renew certificates from Let's Encrypt.

### Option 2: Manual Certificate

1. Obtain certificates from a Certificate Authority (CA)
2. Place the certificate and key files in this directory:
   - `server.crt` - Certificate file
   - `server.key` - Private key file
3. Configure environment variables:

```bash
export TLS_ENABLED=true
export TLS_CERT_PATH=./certs/server.crt
export TLS_KEY_PATH=./certs/server.key
```

## Environment Variables

| Variable | Description | Default |
|----------|-------------|---------|
| `TLS_ENABLED` | Enable HTTPS server | `false` |
| `TLS_CERT_PATH` | Path to TLS certificate | - |
| `TLS_KEY_PATH` | Path to TLS private key | - |
| `TLS_PORT` | HTTPS port | `8443` |
| `HTTPS_ONLY` | Redirect HTTP to HTTPS | `false` |
| `TLS_AUTO_CERT` | Use Let's Encrypt | `false` |
| `TLS_AUTO_CERT_HOST` | Domain for Let's Encrypt | - |

## Security Notes

- Never commit certificate files to version control
- Use strong key sizes (minimum 2048 bits for RSA)
- Regularly rotate certificates
- Monitor certificate expiration dates
- Use HTTPS_ONLY=true in production to force secure connections

## Testing HTTPS

Test the HTTPS connection:

```bash
# Test with curl
curl -k https://localhost:8443/health

# View certificate details
openssl s_client -connect localhost:8443 -showcerts
```