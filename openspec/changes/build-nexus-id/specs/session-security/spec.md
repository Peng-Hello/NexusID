## ADDED Requirements

### Requirement: JWT Signature Algorithm
The system MUST sign all `id_token` and `access_token` JWTs using the RS256 asymmetric encryption algorithm.

#### Scenario: JWT is generated
- **WHEN** a token is issued by the system
- **THEN** its header indicates `alg: RS256` and the signature is valid against the platform's private key

### Requirement: Refresh Token Rotation
The system MUST implement Refresh Token Rotation, issuing a new refresh token and immediately invalidating the old one upon each use.

#### Scenario: Refresh token is consumed
- **WHEN** a user exchanges a valid refresh token
- **THEN** the system issues a new refresh token and marks the previous one as invalid

### Requirement: High Availability Cache Strategy
The system SHALL implement caching mechanisms (like singleflight and negative caching) to prevent cache stampedes and cache penetration.

#### Scenario: High concurrency on absent key
- **WHEN** multiple concurrent requests arrive for a non-existent user ID
- **THEN** the database is queried only once, and a brief NULL cache is set to serve subsequent requests
