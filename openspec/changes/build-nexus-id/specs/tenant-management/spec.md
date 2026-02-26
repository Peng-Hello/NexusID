## ADDED Requirements

### Requirement: Tenant Registration
The system SHALL allow super-administrators to create and provision new tenants within the platform.

#### Scenario: Admin creates a new tenant
- **WHEN** the super admin submits a new tenant request with basic info
- **THEN** the system provisions a unique tenant ID and initializes default configurations

### Requirement: Tenant Data Isolation
The system SHALL ensure that users, clients, and roles belonging to one tenant are not accessible or visible to other tenants.

#### Scenario: Admin queries users
- **WHEN** an admin of Tenant A requests the user list
- **THEN** only users natively belonging to Tenant A are returned

### Requirement: Tenant Client Registration
The system SHALL allow tenant administrators to register OIDC clients (applications) specific to their tenant.

#### Scenario: Admin registers an application
- **WHEN** a tenant administrator creates a new client application
- **THEN** a unique `client_id` and `client_secret` are generated and bound to that tenant
