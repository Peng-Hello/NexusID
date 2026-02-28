<p align="center">
  <img src="./doc/assets/logo.png" alt="NexusID Logo" width="200" />
</p>

# NexusID

*[阅读中文版本 (简体中文)](README_zh.md)*

NexusID is a full-stack Single Sign-On (SSO) platform built with Go (backend) and React (frontend). It provides centralized authentication and identity management using the OpenID Connect (OIDC) protocol.

## Features

### Core Capabilities
- **Internationalization (i18n)**: Full English and Chinese language support
- **OIDC Protocol**: Complete OpenID Connect (OAuth 2.0 +) implementation
- **Multi-Tenant Architecture**: Full data isolation per tenant
- **RS256 JWT Tokens**: Asymmetric encryption with JWKS endpoint
- **Refresh Token Rotation**: Secure token rotation mechanism
- **User Management**: Registration, login, roles, and permissions
- **Anti-Brute-Force**: Rate limiting and account lockout
- **Single Sign-Out**: Global session revocation

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
- React Router for navigation

## Project Structure

```
NexusID/
├── backend/                 # Go backend service
│   ├── internal/
│   │   ├── cache/          # Redis cache wrapper with Singleflight
│   │   ├── config/         # Configuration management
│   │   ├── database/       # Database connection
│   │   ├── dto/            # Data Transfer Objects
│   │   ├── handler/        # HTTP handlers
│   │   ├── jwt/            # JWT token utilities & key management
│   │   ├── logger/         # Zap logger wrapper
│   │   ├── middleware/     # Gin middlewares (rate limiting)
│   │   ├── models/         # GORM models
│   │   ├── repository/     # Database repositories
│   │   └── service/        # Business logic layer
│   ├── keys/               # RSA key pair storage
│   ├── main.go             # Application entry point
│   └── config.yaml         # Configuration file
├── frontend/               # React frontend
│   ├── src/
│   │   ├── components/     # React components (ui/* for shadcn/ui)
│   │   ├── pages/          # Page components
│   │   ├── layouts/        # Layout components
│   │   └── lib/            # Utilities
│   └── package.json
├── deployments/
│   ├── migrations/         # golang-migrate migration files
│   └── mysql/init/         # MySQL initialization scripts
└── docker-compose.yml      # Local development setup
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

2. Start services with Docker Compose:
```bash
docker-compose up -d
```

This will start:
- MySQL on port 3306
- Redis on port 6379
- Backend API on port 8080

3. Run database migrations:
```bash
# Install golang-migrate
go install -tags 'mysql' github.com/golang-migrate/migrate/v4/cmd/migrate@latest

# Run migrations
migrate -path deployments/migrations -database "mysql://nexus_user:nexus_password@tcp(localhost:3306)/nexus_id" up
```

4. Start the frontend:
```bash
cd frontend
npm install
npm run dev
```

The frontend will be available at `http://localhost:5173`

5. Login with the default admin account:

| Field    | Value                |
|----------|----------------------|
| Email    | `admin@nexusid.com`  |
| Password | `admin12345678`      |

> ⚠️ **Important**: Please change the default password immediately after your first login in a production environment.

### Manual Setup

#### Backend

1. Install dependencies:
```bash
cd backend
go mod download
```

2. Configure environment (edit `config.yaml` or set environment variables):
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

3. Run the application:
```bash
go run main.go
```

The backend will be available at `http://localhost:8080`

#### Frontend

1. Install dependencies:
```bash
cd frontend
npm install
```

2. Start development server:
```bash
npm run dev
```

The frontend will be available at `http://localhost:5173`

## API Documentation

### Health Check
- `GET /health` - Health check endpoint
- `GET /ready` - Readiness probe

### OIDC Endpoints
- `GET /.well-known/openid-configuration` - OIDC discovery
- `GET /.well-known/jwks.json` - JWKS public keys
- `GET /oauth/authorize` - Authorization endpoint
- `POST /oauth/token` - Token exchange endpoint
- `POST /oauth/revoke` - Token revocation
- `GET /oauth/logout` - Logout endpoint

### Authentication
- `POST /api/v1/auth/login` - User login
- `POST /api/v1/auth/register` - User registration
- `POST /api/v1/auth/refresh` - Refresh access token

### Tenant Management
- `GET /api/v1/tenants` - List tenants
- `POST /api/v1/tenants` - Create tenant
- `GET /api/v1/tenants/:id` - Get tenant
- `PUT /api/v1/tenants/:id` - Update tenant
- `DELETE /api/v1/tenants/:id` - Delete tenant

### User Management
- `GET /api/v1/tenants/:tenant_id/users` - List users
- `POST /api/v1/tenants/:tenant_id/users` - Create user
- `GET /api/v1/tenants/:tenant_id/users/:id` - Get user
- `PUT /api/v1/tenants/:tenant_id/users/:id` - Update user
- `DELETE /api/v1/tenants/:tenant_id/users/:id` - Delete user

### Role Management
- `GET /api/v1/tenants/:tenant_id/roles` - List roles
- `POST /api/v1/tenants/:tenant_id/roles` - Create role
- `PUT /api/v1/tenants/:tenant_id/roles/:id` - Update role
- `DELETE /api/v1/tenants/:tenant_id/roles/:id` - Delete role
- `POST /api/v1/tenants/:tenant_id/roles/assign` - Assign role to user
- `POST /api/v1/tenants/:tenant_id/roles/revoke` - Revoke role from user

### OIDC Client Management
- `GET /api/v1/tenants/:tenant_id/clients` - List OIDC clients
- `POST /api/v1/tenants/:tenant_id/clients` - Create client
- `PUT /api/v1/tenants/:tenant_id/clients/:id` - Update client
- `DELETE /api/v1/tenants/:tenant_id/clients/:id` - Delete client
- `POST /api/v1/tenants/:tenant_id/clients/:id/rotate-secret` - Rotate client secret

## Configuration

### Backend Configuration (config.yaml)

```yaml
server:
  port: 8080
  mode: debug  # debug | release
  read_timeout: 60
  write_timeout: 60

database:
  host: localhost
  port: 3306
  user: nexus_user
  password: nexus_password
  dbname: nexus_id
  max_open_conns: 100
  max_idle_conns: 10

redis:
  host: localhost
  port: 6379
  password: ""
  db: 0
  pool_size: 100
  min_idle_conn: 10

jwt:
  issuer: http://localhost:8080
  access_token_expiry: 3600     # 1 hour
  refresh_token_expiry: 2592000 # 30 days
  private_key_path: ./keys/private.pem
  public_key_path: ./keys/public.pem
```

### Environment Variables

You can override config values with environment variables (prefixed with `NEXUS_`):

```bash
export NEXUS_SERVER_PORT=8080
export NEXUS_SERVER_MODE=release
export NEXUS_DATABASE_HOST=localhost
export NEXUS_DATABASE_PASSWORD=your_password
# etc.
```

## Security Features

### Rate Limiting
- IP-based rate limiting for login attempts
- Max 5 failed attempts per 15 minutes
- Account lockout for 30 minutes on threshold breach

### Password Security
- Bcrypt hashing for passwords
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

## Production Deployment

### Building for Production

**Backend:**
```bash
cd backend
CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o nexus-id .
```

**Frontend:**
```bash
cd frontend
npm run build
```

### Docker Deployment

Build and run with Docker Compose:
```bash
docker-compose up -d --build
```

### Environment Considerations

For production:
- Set `GIN_MODE=release`
- Use strong database passwords
- Configure proper JWT key paths
- Enable HTTPS/TLS
- Set up Redis clustering/sentinel
- Configure database backups

## License

See LICENSE file for details.
