# IM Backend Service

A real estate platform backend API with property listings, project management, and user authentication. Serves both API endpoints and your React frontend.

## Quick Start

1. **Clone and start the database:**
   ```bash
   git clone <your-repo-url>
   cd IM-backend-GO
   make docker-up
   ```

2. **Set up the database:**
   ```bash
   make migrate      # Creates database tables
   make seed-data    # Adds admin users and test data
   ```

3. **Run the server:**
   ```bash
   make run
   ```

4. **Verify it's working:**
   - Visit `http://localhost:8080/v1/api/health` - should return "OK"
   - API is available at `http://localhost:8080/v1/api/`

That's it! The server is running with a properly set up database.

## Frontend Serving Modes

The backend can serve your React app in different ways:

### 🔄 Development Mode (Recommended for development)
Proxies to your React dev server with hot reload:

```bash
# Start your React app (usually npm start)
FRONTEND_PROXY_URL=http://localhost:3000 make run
```

### ⚡ Production Mode (Fast serving from memory)
Downloads and serves your app from a zip file:

```bash
STATIC_ASSETS_URL=https://example.com/your-app.zip make run
```

### 📁 Basic Mode (Serves from build folder)
Serves files from the `build/` directory:

```bash
make run  # Default mode, no config needed. Another way to say, API Only
```

## Configuration

Create a `.env` file for custom settings:

```env
# Database (optional - uses Docker by default)
DATABASE_URL=root:password@tcp(localhost:3306)/mydb

# Server
PORT=8080

# Authentication (required for login features)
AUTH_JWT_SECRET=your-secret-key-here

# AWS S3 (required for file uploads)
AWS_BUCKET=your-bucket-name
AWS_ACCESS_KEY_ID=your-access-key
AWS_SECRET_KEY=your-secret-key

# Frontend serving (optional)
FRONTEND_PROXY_URL=http://localhost:3000
STATIC_ASSETS_URL=https://example.com/app.zip
```

## Common Commands

```bash
# Start the server
make run

# Database setup
make migrate       # Create database tables
make seed-data     # Add admin users (admin/admin123, business_partner/bp123, etc.)
make seed-projects # Add sample properties and projects

# Reset database (clean start with all data)
make reset

# Stop database
make docker-down
```

## Test Users (After running seed-data)

The seeding creates these test users for development:

```
Superadmin:
- Username: admin
- Password: admin123

Business Partner:
- Username: business_partner  
- Password: bp123

Digital Marketing:
- Username: digital_marketing
- Password: dm123
```

Use these credentials to test authentication and different user roles.

## Troubleshooting

**Server won't start:**
- Make sure Docker is running
- Run `make docker-up` to start the database
- Check if port 8080 is available

**Frontend not loading:**
- For development: make sure your React app is running on localhost:3000
- For production: verify your zip URL is accessible
- For basic mode: you'll get 404 errors if no `build/` directory exists

**Database connection issues:**
- Run `make reset` to recreate the database with fresh data
- Make sure you ran `make migrate` after starting the database
- Check your `DATABASE_URL` if using custom database

**Empty database or missing tables:**
- Run `make migrate` to create database tables
- Run `make seed-data` to add test users
- Run `make seed-projects` to add sample data

## Need Help?

- Check the logs when running `make run`
- Verify your `.env` file configuration

---

*For detailed technical documentation, see CLAUDE.md*