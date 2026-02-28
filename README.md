<p align="center">
  <img src="./doc/assets/logo.png" alt="NexusID Logo" width="200" />
</p>

# NexusID

*[中文版 (简体中文)](README_zh.md)* ｜ *[📖 User Guide](doc/user_guide.md)* ｜ *[🔧 Developer Guide](doc/developer_guide.md)*

NexusID is a full-stack Single Sign-On (SSO) platform built with Go (backend) and React (frontend). It provides centralized authentication and identity management using the OpenID Connect (OIDC) protocol.

## Features

### Core Capabilities
- **Internationalization (i18n)**: Full English and Chinese language support with seamless switching
- **OIDC Protocol**: Complete OpenID Connect (OAuth 2.0 +) implementation
- **Multi-Tenant Architecture**: Full data isolation per tenant
- **RS256 JWT Tokens**: Asymmetric encryption mechanism with JWKS endpoint
- **Refresh Token Rotation**: Secure refresh token rotation and expiration mechanism
- **User Management**: Registration, login, roles, and permissions assignment
- **Anti-Brute-Force**: Login rate limiting strategy and automatic account lockout
- **Single Sign-Out**: Global session logout and Token revocation support
- **Password Change**: User self-service password change functionality

### Tech Stack

**Backend:**
- Go 1.23+ with Gin web framework
- MySQL 8.0 (no foreign keys for performance)
- Redis 7.0 with Singleflight cache wrapper
- GORM for database ORM
- Zap for structured logging
- Viper for configuration management
- RS256 JWT tokens

**Frontend:**
- React 18 with TypeScript
- Vite for build tooling
- Tailwind CSS for styling
- shadcn/ui component library
- React Router for page routing
- react-i18next for complete internationalization support

## Project Structure

```
NexusID/
├── backend/                 # Go backend service
│   ├── internal/
│   │   ├── cache/          # Redis cache wrapper with Singleflight
│   │   ├── config/         # Configuration file mapping management
│   │   ├── database/       # Database connection
│   │   ├── dto/            # Data Transfer Objects
│   │   ├── handler/        # HTTP route controllers
│   │   ├── jwt/            # JWT generation and key management
│   │   ├── logger/         # Zap logger wrapper
│   │   ├── middleware/     # Gin middlewares (including rate limiting and auth)
│   │   ├── models/         # GORM entity models
│   │   ├── repository/     # Database repository layer
│   │   └── service/        # Core business logic layer
│   ├── keys/               # RSA key pair storage directory
│   ├── main.go             # Application entry point
│   └── config.yaml         # Server configuration file
├── frontend/               # React frontend
│   ├── src/
│   │   ├── components/     # React base components (including shadcn/ui)
│   │   ├── i18n/           # Internationalization config and language packs
│   │   ├── pages/          # View-level page components
│   │   ├── layouts/        # Common layout framework
│   │   └── lib/            # Utility functions and configuration
│   └── package.json
├── deployments/
│   ├── migrations/         # golang-migrate migration scripts
│   └── mysql/init/         # MySQL database initialization scripts
└── docker-compose.yml      # One-click local environment setup
```

## Quick Start

### Prerequisites
- Go 1.23+
- Node.js 22+
- MySQL 8.0+
- Redis 7.0+
- Docker & Docker Compose (optional)

### Using Docker Compose (Recommended)

1. Clone and navigate to the project:
```bash
cd NexusID
```

2. Start all services:
```bash
docker-compose up -d
```

This will start the following services:
- MySQL running on port 3306
- Redis running on port 6379
- Backend API running on port 8080

3. Run database migrations:
```bash
# Install golang-migrate
go install -tags 'mysql' github.com/golang-migrate/migrate/v4/cmd/migrate@latest

# Run migration scripts
migrate -path deployments/migrations -database "mysql://nexus_user:nexus_password@tcp(localhost:3306)/nexus_id" up
```

4. Start the frontend:
```bash
cd frontend
npm install
npm run dev
```

The frontend service will start at `http://localhost:5173`

5. Login with the default admin account:

| Field    | Value                |
|----------|----------------------|
| Email    | `admin@nexusid.com`  |
| Password | `admin12345678`      |

> ⚠️ **Important**: Please change the default password immediately after your first login in a production environment.

### Manual Deployment Guide

#### Backend

1. Install dependencies:
```bash
cd backend
go mod download
```

2. Configure environment variables (edit `config.yaml` or set via system environment variables):
```yaml
server:
  port: 8080
  mode: debug

database:
  host: localhost
  port: 3306
  user: nexus_user
  password: nexus_password
  dbname: nexus_id

redis:
  host: localhost
  port: 6379

jwt:
  issuer: http://localhost:8080
  private_key_path: ./keys/private.pem
  public_key_path: ./keys/public.pem
```

3. Start the backend server:
```bash
go run main.go
```

The backend API will be available at `http://localhost:8080`

#### Frontend

1. Install package dependencies:
```bash
cd frontend
npm install
```

2. Start development server:
```bash
npm run dev
```

The frontend interface will be available at `http://localhost:5173`

## API Documentation Summary

### Health Check
- `GET /health` - Liveness probe
- `GET /ready` - Readiness probe

### OIDC Core Endpoints
- `GET /.well-known/openid-configuration` - OIDC service discovery
- `GET /.well-known/jwks.json` - JWKS public key distribution
- `GET /oauth/authorize` - Authorization endpoint
- `POST /oauth/token` - Token exchange/refresh endpoint
- `POST /oauth/revoke` - Token revocation endpoint
- `GET /oauth/logout` - Single sign-out endpoint

### Authentication
- `POST /api/v1/auth/login` - User login
- `POST /api/v1/auth/register` - User registration
- `POST /api/v1/auth/refresh` - Manual access token refresh
- `POST /api/v1/user/change-password` - Change user password

### Other Management APIs
*(Refer to the detailed API documentation for complete CRUD operations on tenants, users, roles, and OIDC clients)*

## Security Features

### Rate Limiting
- IP-based rate limiting for login attempts
- Maximum 5 failed attempts per 15 minutes
- Account lockout for 30 minutes when threshold is reached

### Password Security
- Bcrypt hashing for password storage
- Minimum 8 character password requirement
- Failed login attempt tracking

### Token Security
- RS256 asymmetric encryption (no shared secrets)
- JWKS endpoint for public key distribution
- Refresh token rotation (old tokens invalidated on use)
- Token expiration and revocation

### Tenant Isolation
- All queries scoped to tenant_id
- Service layer validation
- Database-level constraints

## Development

### Running Tests

```bash
# Backend tests
cd backend
go test ./...

# Frontend tests
cd frontend
npm test
```

### Database Migrations

Create new migration:
```bash
migrate create -ext sql -dir deployments/migrations -seq migration_name
```

Run migrations:
```bash
migrate -path deployments/migrations -database "mysql://user:pass@tcp(localhost:3306)/nexus_id" up
```

Rollback migrations:
```bash
migrate -path deployments/migrations -database "mysql://user:pass@tcp(localhost:3306)/nexus_id" down 1
```

## Production Deployment Recommendations

### Build Artifacts

**Backend binary:**
```bash
cd backend
CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o nexus-id .
```

**Frontend static files:**
```bash
cd frontend
npm run build
```

### Environment Considerations

Production checklist:
- Set `GIN_MODE` to `release`
- Enable HTTPS/TLS
- Update to strong database passwords
- Modify `config.yaml` to disable CORS allowing all origins
- Configure proper RS256 security key mounting
- Establish database backup mechanism

## License

See LICENSE file for details.
