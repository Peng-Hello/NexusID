package service

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"time"

	"github.com/nexus-id/backend/internal/jwt"
	"github.com/nexus-id/backend/internal/models"
	"github.com/nexus-id/backend/internal/repository"
	"gorm.io/gorm"
)

var (
	ErrInvalidClient      = errors.New("invalid client")
	ErrInvalidRedirectURI = errors.New("invalid redirect URI")
	ErrInvalidCode        = errors.New("invalid authorization code")
	ErrInvalidGrant       = errors.New("invalid grant")
	ErrExpiredCode        = errors.New("authorization code expired")
)

// OIDCService handles OIDC protocol operations
type OIDCService struct {
	authCodeRepo     *repository.AuthorizationCodeRepository
	refreshTokenRepo *repository.RefreshTokenRepository
	sessionRepo      *repository.SessionRepository
	userRepo         *repository.UserRepository
	tenantRepo       *repository.TenantRepository
	clientRepo       *repository.OIDCClientRepository
	roleRepo         *repository.RoleRepository
	tokenService     *jwt.TokenService
}

// NewOIDCService creates a new OIDCService
func NewOIDCService(tokenService *jwt.TokenService) *OIDCService {
	return &OIDCService{
		authCodeRepo:     repository.NewAuthorizationCodeRepository(),
		refreshTokenRepo: repository.NewRefreshTokenRepository(),
		sessionRepo:      repository.NewSessionRepository(),
		userRepo:         repository.NewUserRepository(),
		tenantRepo:       repository.NewTenantRepository(),
		clientRepo:       repository.NewOIDCClientRepository(),
		roleRepo:         repository.NewRoleRepository(),
		tokenService:     tokenService,
	}
}

// AuthorizeRequest represents an authorization request
type AuthorizeRequest struct {
	ResponseType        string `form:"response_type" binding:"required"`
	ClientID            string `form:"client_id" binding:"required"`
	RedirectURI         string `form:"redirect_uri" binding:"required"`
	Scope               string `form:"scope"`
	State               string `form:"state"`
	CodeChallenge       string `form:"code_challenge"`
	CodeChallengeMethod string `form:"code_challenge_method"`
	TenantID            int64
	UserID              int64
}

// TokenRequest represents a token request
type TokenRequest struct {
	GrantType    string `form:"grant_type" binding:"required"`
	Code         string `form:"code"`
	RedirectURI  string `form:"redirect_uri"`
	ClientID     string `form:"client_id"`
	ClientSecret string `form:"client_secret"`
	RefreshToken string `form:"refresh_token"`
	CodeVerifier string `form:"code_verifier"`
}

// TokenResponse represents a token response
type TokenResponse struct {
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int64  `json:"expires_in"`
	RefreshToken string `json:"refresh_token,omitempty"`
	IDToken      string `json:"id_token"`
}

// GenerateAuthorizationCode generates an authorization code for the OAuth2 flow
func (s *OIDCService) GenerateAuthorizationCode(req *AuthorizeRequest) (string, error) {
	// Validate response type
	if req.ResponseType != "code" {
		return "", errors.New("unsupported response type")
	}

	// Get client
	client, err := s.clientRepo.GetByClientID(req.ClientID)
	if err != nil {
		return "", ErrInvalidClient
	}

	if !client.IsActive {
		return "", ErrInvalidClient
	}

	// Validate redirect URI
	if !s.isValidRedirectURI(client.RedirectURIs, req.RedirectURI) {
		return "", ErrInvalidRedirectURI
	}

	// Validate PKCE if code challenge is present
	if req.CodeChallenge != "" && req.CodeChallengeMethod != "plain" && req.CodeChallengeMethod != "S256" {
		return "", errors.New("invalid code challenge method")
	}

	// Generate authorization code
	code := s.authCodeRepo.GenerateCode()

	authCode := &models.AuthorizationCode{
		TenantID:            client.TenantID,
		UserID:              req.UserID,
		ClientID:            req.ClientID,
		Code:                code,
		RedirectURI:         req.RedirectURI,
		Scope:               req.Scope,
		CodeChallenge:       req.CodeChallenge,
		CodeChallengeMethod: req.CodeChallengeMethod,
		ExpiresAt:           time.Now().UTC().Add(10 * time.Minute),
		IsUsed:              false,
	}

	if err := s.authCodeRepo.Create(authCode); err != nil {
		return "", fmt.Errorf("failed to create authorization code: %w", err)
	}

	return code, nil
}

// ExchangeCodeForToken exchanges an authorization code for tokens
func (s *OIDCService) ExchangeCodeForToken(req *TokenRequest) (*TokenResponse, error) {
	// Validate grant type
	if req.GrantType != "authorization_code" {
		return nil, ErrInvalidGrant
	}

	// Get and validate client
	client, err := s.clientRepo.GetByClientID(req.ClientID)
	if err != nil {
		return nil, ErrInvalidClient
	}

	// Validate client secret (skip for public clients)
	if !client.IsPublic {
		if err := s.validateClientSecret(client, req.ClientSecret); err != nil {
			return nil, ErrInvalidClient
		}
	}

	// Get authorization code
	authCode, err := s.authCodeRepo.GetByCode(req.Code)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrInvalidCode
		}
		return nil, err
	}

	// Validate code hasn't expired
	if authCode.IsExpired() {
		return nil, ErrExpiredCode
	}

	// Validate code hasn't been used
	if authCode.IsUsed {
		return nil, ErrInvalidCode
	}

	// Validate redirect URI matches
	if authCode.RedirectURI != req.RedirectURI {
		return nil, ErrInvalidRedirectURI
	}

	// Validate client ID matches
	if authCode.ClientID != req.ClientID {
		return nil, ErrInvalidClient
	}

	// Validate PKCE code verifier if challenge was present
	if authCode.CodeChallenge != "" {
		if err := s.validatePKCE(authCode.CodeChallenge, authCode.CodeChallengeMethod, req.CodeVerifier); err != nil {
			return nil, err
		}
	}

	// Mark code as used
	if err := s.authCodeRepo.MarkAsUsed(req.Code); err != nil {
		return nil, err
	}

	// Get user with roles
	user, err := s.userRepo.GetByID(authCode.UserID)
	if err != nil {
		return nil, err
	}

	// Get user roles
	roles := make([]string, len(user.Roles))
	for i, role := range user.Roles {
		roles[i] = role.Name
	}

	// Generate tokens
	accessToken, err := s.tokenService.GenerateAccessToken(user, req.ClientID, roles)
	if err != nil {
		return nil, err
	}

	idToken, err := s.tokenService.GenerateIDToken(user, req.ClientID, "", roles)
	if err != nil {
		return nil, err
	}

	// Generate and store refresh token
	refreshTokenStr, err := s.generateAndStoreRefreshToken(user, req.ClientID)
	if err != nil {
		return nil, err
	}

	// Create session
	_ = s.createSession(user, req.ClientID)

	return &TokenResponse{
		AccessToken:  accessToken,
		TokenType:    "Bearer",
		ExpiresIn:    int64(s.tokenService.GetTokenExpiry().Seconds()),
		RefreshToken: refreshTokenStr,
		IDToken:      idToken,
	}, nil
}

// RefreshAccessToken refreshes an access token using a refresh token
func (s *OIDCService) RefreshAccessToken(req *TokenRequest) (*TokenResponse, error) {
	// Validate grant type
	if req.GrantType != "refresh_token" {
		return nil, ErrInvalidGrant
	}

	// Get and validate client
	client, err := s.clientRepo.GetByClientID(req.ClientID)
	if err != nil {
		return nil, ErrInvalidClient
	}

	// Validate client secret (skip for public clients)
	if !client.IsPublic {
		if err := s.validateClientSecret(client, req.ClientSecret); err != nil {
			return nil, ErrInvalidClient
		}
	}

	// Get refresh token from storage
	refreshToken, err := s.refreshTokenRepo.GetByTokenHash(req.RefreshToken)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrInvalidGrant
		}
		return nil, err
	}

	// Check if revoked
	if refreshToken.IsRevoked {
		return nil, ErrInvalidGrant
	}

	// Check if expired
	if refreshToken.IsExpired() {
		return nil, ErrInvalidGrant
	}

	// Validate client ID matches
	if refreshToken.ClientID != req.ClientID {
		return nil, ErrInvalidClient
	}

	// Get user with roles
	user, err := s.userRepo.GetByID(refreshToken.UserID)
	if err != nil {
		return nil, err
	}

	// Get user roles
	roles := make([]string, len(user.Roles))
	for i, role := range user.Roles {
		roles[i] = role.Name
	}

	// Revoke old refresh token (token rotation)
	if err := s.refreshTokenRepo.Revoke(req.RefreshToken); err != nil {
		return nil, err
	}

	// Revoke all other refresh tokens for this user (security)
	s.refreshTokenRepo.RevokeAllForUserExcept(refreshToken.UserID, req.RefreshToken)

	// Revoke all other sessions for this user
	s.sessionRepo.RevokeAllForUserExcept(refreshToken.UserID, "")

	// Generate new tokens
	accessToken, err := s.tokenService.GenerateAccessToken(user, req.ClientID, roles)
	if err != nil {
		return nil, err
	}

	idToken, err := s.tokenService.GenerateIDToken(user, req.ClientID, "", roles)
	if err != nil {
		return nil, err
	}

	// Generate and store new refresh token
	newRefreshTokenStr, err := s.generateAndStoreRefreshToken(user, req.ClientID)
	if err != nil {
		return nil, err
	}

	return &TokenResponse{
		AccessToken:  accessToken,
		TokenType:    "Bearer",
		ExpiresIn:    int64(s.tokenService.GetTokenExpiry().Seconds()),
		RefreshToken: newRefreshTokenStr,
		IDToken:      idToken,
	}, nil
}

// RevokeToken handles RFC 7009 token revocation
func (s *OIDCService) RevokeToken(token string, tokenTypeHint string) error {
	// Typically, we only care about revoking the refresh token since
	// access tokens in our implementation are stateless JWTs.
	// If hint is access_token, we can't really revoke it securely without a denylist.
	// We will attempt to revoke it as a refresh token from the database.

	err := s.refreshTokenRepo.Revoke(token)
	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		// RFC says to return 200 OK even if the token doesn't exist
		return nil
	}
	return err
}

// RevokeSession revokes a user session (single sign-out)
func (s *OIDCService) RevokeSession(sessionID string) error {
	return s.sessionRepo.Revoke(sessionID)
}

// RevokeAllUserSessions revokes all sessions for a user
func (s *OIDCService) RevokeAllUserSessions(userID int64) error {
	return s.sessionRepo.RevokeAllForUser(userID)
}

// Helper methods

func (s *OIDCService) isValidRedirectURI(allowedURIs []string, providedURI string) bool {
	for _, uri := range allowedURIs {
		if uri == providedURI {
			return true
		}
	}
	return false
}

func (s *OIDCService) validateClientSecret(client *models.OIDCClient, secret string) error {
	// This would use bcrypt to validate the secret
	// For now, return nil (should be implemented properly)
	return nil
}

func (s *OIDCService) validatePKCE(codeChallenge, method, verifier string) error {
	// Validate PKCE code verifier
	if method == "plain" {
		if codeChallenge != verifier {
			return errors.New("invalid code verifier")
		}
	} else if method == "S256" {
		// SHA256(verifier) should match codeChallenge
		// Implementation needed
	}
	return nil
}

func (s *OIDCService) generateAndStoreRefreshToken(user *models.User, clientID string) (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	token := base64.URLEncoding.EncodeToString(b)

	refreshToken := &models.RefreshToken{
		TenantID:  user.TenantID,
		UserID:    user.ID,
		ClientID:  clientID,
		TokenHash: token,
		IsRevoked: false,
		ExpiresAt: time.Now().UTC().Add(30 * 24 * time.Hour), // 30 days
	}

	if err := s.refreshTokenRepo.Create(refreshToken); err != nil {
		return "", err
	}

	return token, nil
}

func (s *OIDCService) createSession(user *models.User, clientID string) string {
	sessionID := s.sessionRepo.GenerateSessionID()

	session := &models.Session{
		TenantID:  user.TenantID,
		UserID:    user.ID,
		ClientID:  clientID,
		SessionID: sessionID,
		IsRevoked: false,
		ExpiresAt: time.Now().UTC().Add(24 * time.Hour), // 24 hours
	}

	s.sessionRepo.Create(session)

	return sessionID
}
