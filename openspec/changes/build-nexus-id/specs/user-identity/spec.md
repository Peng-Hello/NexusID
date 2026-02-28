## ADDED Requirements

### Requirement: User Registration and Login
The system SHALL allow users to register and authenticate using an email/password combination under their specific tenant.

#### Scenario: Existing user logs in
- **WHEN** a user provides valid credentials
- **THEN** they get authenticated and a secure session is established

### Requirement: Identity Roles Assignment
The system SHALL allow administrators to define macro-roles and assign them to users.

#### Scenario: User receives a role
- **WHEN** an admin assigns the "Internal Employee" role to a user
- **THEN** subsequent ID Tokens for that user include the "Internal Employee" role claim

### Requirement: Account Security Measures
The system SHALL enforce password minimum length and temporary account lockout after multiple failed login attempts.

#### Scenario: Multiple failed login attempts
- **WHEN** a user enters incorrect passwords consecutive times exceeding the threshold
- **THEN** the account is locked for the predefined lockout duration
