# IM Backend Service

A modern Go backend service with JWT authentication, OTP-based admin login, and MySQL database integration.

## Features

- 🚀 Built with Go and Gorilla Mux
- 🔐 JWT-based authentication
- 📱 Phone OTP for admin login
- 🗄️ MySQL database with Ent ORM
- 🧪 Comprehensive test suite
- 🐳 Docker support for local development
- 🔄 Database migrations with Atlas

## Prerequisites

- Go 1.21 or later
- Docker and Docker Compose
- Make (optional, but recommended)

## Quick Start

1. Clone the repository:
```bash
git clone https://github.com/yourusername/im_backend_go.git
cd im_backend_go
```

2. Start the MySQL database:
```bash
make docker-up
```

3. Run database migrations:
```bash
make migrate
```

4. Start the server:
```bash
make run
```

The server will start on `http://localhost:8080`.

## Development

### Project Structure

```
.
├── cmd/
│   └── server/          # Application entry point
├── internal/
│   ├── auth/           # Authentication logic
│   ├── config/         # Configuration management
│   ├── database/       # Database connection and setup
│   ├── handlers/       # HTTP request handlers
│   ├── middleware/     # HTTP middleware
│   ├── router/         # Route definitions
│   └── testutil/       # Test utilities
├── ent/                # Ent ORM schema and generated code
├── migrations/         # Database migrations
└── Makefile           # Build and development commands
```

### Available Commands

```bash
# Start the application
make run

# Run tests
make test

# Build the application
make build

# Start Docker services
make docker-up

# Stop Docker services
make docker-down

# Run database migrations
make migrate

# Clean build artifacts
make clean
```

### Environment Variables

Create a `.env` file in the project root:

```env
DB_DSN=root:password@tcp(localhost:3306)/im_db
PORT=8080
JWT_SECRET=your-secret-key
```

## API Endpoints

### Public Endpoints

- `GET /health` - Health check endpoint
- `POST /auth/token` - Generate JWT token (requires phone and OTP)

### Protected Endpoints (Admin Only)

- `POST /admin/users` - Create a new user
- `GET /admin/users` - List all users
- `GET /admin/users/{id}` - Get user details
- `PUT /admin/users/{id}` - Update user
- `DELETE /admin/users/{id}` - Delete user

## Testing

Run the test suite:

```bash
make test
```

The test suite includes:
- Integration tests with a test database
- API endpoint tests
- Authentication tests

## Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details. 