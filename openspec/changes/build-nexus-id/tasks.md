## 1. Project Core & Infrastructure Setup

- [x] 1.1 Initialize full-stack monorepo structure (backend/ and frontend/ dirs)
- [x] 1.2 Initialize Go module and configure basic HTTP routing (Gin)
- [x] 1.3 Configure Docker Compose with MySQL 8.0 and Redis 7.0
- [x] 1.4 Implement structured logging (zap) and configuration management (Viper)

## 2. Database & Data Models

- [x] 2.1 Design database schema and create migration scripts (golang-migrate)
- [x] 2.2 Setup GORM connection and implement repositories with decapsulated DTOs
- [x] 2.3 Implement Redis connection pool and High Availability Cache Wrapper (Singleflight & NULL caching)

## 3. OIDC Core & Security Foundation

- [x] 3.1 Implement RS256 Key Pair generation, loading, and secure storage
- [x] 3.2 Expose JWKS endpoint (`/.well-known/jwks.json`)
- [x] 3.3 Expose OIDC Discovery endpoint (`/.well-known/openid-configuration`)
- [x] 3.4 Implement JWT encoding/decoding utilities (access_token, id_token)

## 4. Tenant & Client Management

- [x] 4.1 Implement CRUD APIs for Tenant Management
- [x] 4.2 Implement CRUD APIs for OIDC Client Registration
- [x] 4.3 Ensure downstream database query isolation by tenant_id

## 5. User & Identity System

- [x] 5.1 Implement User Registration and Login API (with password hashing)
- [x] 5.2 Implement basic Identity Roles creation API
- [x] 5.3 Implement User to Role assignment API
- [x] 5.4 Implement anti-brute-force rate limiting and lockout mechanism

## 6. OIDC Flows & Session Management

- [x] 6.1 Implement Authorization Code generation and Validation
- [x] 6.2 Implement Token Exchange Endpoint (`/token`)
- [x] 6.3 Implement Refresh Token Rotation mechanism
- [x] 6.4 Implement Single Sign-Out / Session Revocation Endpoint

## 7. Frontend Administration

- [x] 7.1 Initialize React frontend with Tailwind CSS and shadcn/ui
- [x] 7.2 Build unified Login, Register, and Forgot Password views
- [x] 7.3 Build Super Admin Dashboard for Tenant & Client management
- [x] 7.4 Build User & Role Management interface for Tenant Admins
