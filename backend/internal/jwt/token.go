package jwt

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/nexus-id/backend/internal/models"
)

// TokenType represents the type of JWT token
type TokenType string

const (
	AccessTokenType  TokenType = "access_token"
	IDTokenType      TokenType = "id_token"
	RefreshTokenType TokenType = "refresh_token"
)

// TokenClaims represents the standard JWT claims
type TokenClaims struct {
	jwt.RegisteredClaims

	// Custom claims
	TenantID  int64     `json:"tenant_id"`
	UserID    int64     `json:"user_id"`
	Email     string    `json:"email"`
	FullName  string    `json:"name,omitempty"`
	Roles     []string  `json:"roles,omitempty"`
	ClientID  string    `json:"client_id,omitempty"`
	TokenType TokenType `json:"token_type"`
	Nonce     string    `json:"nonce,omitempty"` // OIDC nonce
}

// TokenService handles JWT token generation and validation
type TokenService struct {
	keyManager    *KeyManager
	issuer        string
	accessExpiry  time.Duration
	refreshExpiry time.Duration
}

// NewTokenService creates a new TokenService
func NewTokenService(keyManager *KeyManager, issuer string, accessExpiry, refreshExpiry int) *TokenService {
	return &TokenService{
		keyManager:    keyManager,
		issuer:        issuer,
		accessExpiry:  time.Duration(accessExpiry) * time.Second,
		refreshExpiry: time.Duration(refreshExpiry) * time.Second,
	}
}

// GenerateAccessToken generates an access token
func (s *TokenService) GenerateAccessToken(user *models.User, clientID string, roles []string) (string, error) {
	now := time.Now()
	claims := TokenClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    s.issuer,
			Subject:   fmt.Sprintf("%d", user.ID),
			Audience:  []string{clientID},
			ExpiresAt: jwt.NewNumericDate(now.Add(s.accessExpiry)),
			NotBefore: jwt.NewNumericDate(now),
			IssuedAt:  jwt.NewNumericDate(now),
		},
		TenantID:  user.TenantID,
		UserID:    user.ID,
		Email:     user.Email,
		FullName:  user.FullName,
		Roles:     roles,
		ClientID:  clientID,
		TokenType: AccessTokenType,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	tokenString, err := token.SignedString(s.keyManager.GetPrivateKey())
	if err != nil {
		return "", fmt.Errorf("failed to sign access token: %w", err)
	}

	return tokenString, nil
}

// GenerateIDToken generates an ID token (OIDC)
func (s *TokenService) GenerateIDToken(user *models.User, clientID string, nonce string, roles []string) (string, error) {
	now := time.Now()
	claims := TokenClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    s.issuer,
			Subject:   fmt.Sprintf("%d", user.ID),
			Audience:  []string{clientID},
			ExpiresAt: jwt.NewNumericDate(now.Add(s.accessExpiry)),
			NotBefore: jwt.NewNumericDate(now),
			IssuedAt:  jwt.NewNumericDate(now),
		},
		TenantID:  user.TenantID,
		UserID:    user.ID,
		Email:     user.Email,
		FullName:  user.FullName,
		Roles:     roles,
		ClientID:  clientID,
		TokenType: IDTokenType,
		Nonce:     nonce,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	tokenString, err := token.SignedString(s.keyManager.GetPrivateKey())
	if err != nil {
		return "", fmt.Errorf("failed to sign ID token: %w", err)
	}

	return tokenString, nil
}

// GenerateRefreshToken generates a refresh token
func (s *TokenService) GenerateRefreshToken(user *models.User, clientID string) (string, error) {
	now := time.Now()
	claims := TokenClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    s.issuer,
			Subject:   fmt.Sprintf("%d", user.ID),
			Audience:  []string{clientID},
			ExpiresAt: jwt.NewNumericDate(now.Add(s.refreshExpiry)),
			NotBefore: jwt.NewNumericDate(now),
			IssuedAt:  jwt.NewNumericDate(now),
		},
		TenantID:  user.TenantID,
		UserID:    user.ID,
		ClientID:  clientID,
		TokenType: RefreshTokenType,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	tokenString, err := token.SignedString(s.keyManager.GetPrivateKey())
	if err != nil {
		return "", fmt.Errorf("failed to sign refresh token: %w", err)
	}

	return tokenString, nil
}

// ValidateToken validates a JWT token and returns the claims
func (s *TokenService) ValidateToken(tokenString string) (*TokenClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &TokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		// Validate signing method
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return s.keyManager.GetPublicKey(), nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to parse token: %w", err)
	}

	claims, ok := token.Claims.(*TokenClaims)
	if !ok || !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	return claims, nil
}

// GetTokenExpiry returns the expiry duration for access tokens
func (s *TokenService) GetTokenExpiry() time.Duration {
	return s.accessExpiry
}

// GetRefreshTokenExpiry returns the expiry duration for refresh tokens
func (s *TokenService) GetRefreshTokenExpiry() time.Duration {
	return s.refreshExpiry
}
