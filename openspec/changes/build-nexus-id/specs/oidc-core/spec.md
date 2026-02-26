## ADDED Requirements

### Requirement: OIDC Discovery Endpoint
The system SHALL expose an OIDC discovery endpoint (`/.well-known/openid-configuration`) containing provider metadata including the JWKS URI.

#### Scenario: Client requests provider configuration
- **WHEN** a client makes a GET request to `/.well-known/openid-configuration`
- **THEN** it receives a JSON response with OIDC metadata including `jwks_uri`

### Requirement: JWKS Endpoint
The system SHALL expose a JWKS endpoint (`/.well-known/jwks.json`) containing the public keys used to verify RS256 signatures of issued JWTs.

#### Scenario: Client fetches public keys
- **WHEN** a client makes a GET request to `/.well-known/jwks.json`
- **THEN** it receives a JSON response with a list of currently active public keys

### Requirement: Authorization Code Grant
The system SHALL support the OAuth 2.0 Authorization Code flow for authenticating users and issuing tokens.

#### Scenario: Successful authorization code exchange
- **WHEN** a client exchanges a valid authorization code at the `/token` endpoint
- **THEN** the system issues an `access_token`, an `id_token`, and a `refresh_token`

### Requirement: Single Sign-Out
The system SHALL support logging out users globally across registered clients using Back-Channel Logout mechanics or a centralized logout endpoint.

#### Scenario: User initiates logout
- **WHEN** a user logs out from the IdP
- **THEN** their session is terminated and active refresh tokens are revoked
