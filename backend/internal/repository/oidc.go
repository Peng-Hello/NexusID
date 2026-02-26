package repository

import (
	"crypto/sha256"
	"encoding/base64"
	"time"

	"github.com/nexus-id/backend/internal/models"
)

// AuthorizationCodeRepository handles authorization code data operations
type AuthorizationCodeRepository struct {
	*BaseRepository
}

// NewAuthorizationCodeRepository creates a new AuthorizationCodeRepository
func NewAuthorizationCodeRepository() *AuthorizationCodeRepository {
	return &AuthorizationCodeRepository{
		BaseRepository: NewBaseRepository(),
	}
}

// GenerateCode generates a secure authorization code
func (r *AuthorizationCodeRepository) GenerateCode() string {
	b := make([]byte, 32)
	r.db.Raw("SELECT RANDOMBLOB(32)").Scan(&b)
	return base64.URLEncoding.EncodeToString(b)
}

// HashCode hashes an authorization code for storage
func (r *AuthorizationCodeRepository) HashCode(code string) string {
	hash := sha256.Sum256([]byte(code))
	return base64.URLEncoding.EncodeToString(hash[:])
}

// Create creates a new authorization code
func (r *AuthorizationCodeRepository) Create(authCode *models.AuthorizationCode) error {
	authCode.Code = r.HashCode(authCode.Code)
	return r.db.Create(authCode).Error
}

// GetByCode retrieves an authorization code by code
func (r *AuthorizationCodeRepository) GetByCode(code string) (*models.AuthorizationCode, error) {
	hashedCode := r.HashCode(code)
	var authCode models.AuthorizationCode
	err := r.db.Where("code = ? AND is_used = false AND expires_at > ?", hashedCode, time.Now().UTC()).
		Preload("User").
		Preload("Tenant").
		First(&authCode).Error
	if err != nil {
		return nil, err
	}
	return &authCode, nil
}

// MarkAsUsed marks an authorization code as used
func (r *AuthorizationCodeRepository) MarkAsUsed(code string) error {
	hashedCode := r.HashCode(code)
	return r.db.Model(&models.AuthorizationCode{}).
		Where("code = ?", hashedCode).
		Update("is_used", true).Error
}

// DeleteExpired deletes expired authorization codes
func (r *AuthorizationCodeRepository) DeleteExpired() error {
	return r.db.Where("expires_at < ? OR is_used = true", time.Now().UTC()).
		Delete(&models.AuthorizationCode{}).Error
}

// RefreshTokenRepository handles refresh token data operations
type RefreshTokenRepository struct {
	*BaseRepository
}

// NewRefreshTokenRepository creates a new RefreshTokenRepository
func NewRefreshTokenRepository() *RefreshTokenRepository {
	return &RefreshTokenRepository{
		BaseRepository: NewBaseRepository(),
	}
}

// HashToken hashes a refresh token for storage
func (r *RefreshTokenRepository) HashToken(token string) string {
	hash := sha256.Sum256([]byte(token))
	return base64.URLEncoding.EncodeToString(hash[:])
}

// Create creates a new refresh token
func (r *RefreshTokenRepository) Create(refreshToken *models.RefreshToken) error {
	refreshToken.TokenHash = r.HashToken(refreshToken.TokenHash)
	return r.db.Create(refreshToken).Error
}

// GetByTokenHash retrieves a refresh token by hash
func (r *RefreshTokenRepository) GetByTokenHash(tokenHash string) (*models.RefreshToken, error) {
	hashedToken := r.HashToken(tokenHash)
	var refreshToken models.RefreshToken
	err := r.db.Where("token_hash = ? AND is_revoked = false", hashedToken).
		Preload("User").
		Preload("Tenant").
		First(&refreshToken).Error
	if err != nil {
		return nil, err
	}
	return &refreshToken, nil
}

// Revoke revokes a refresh token
func (r *RefreshTokenRepository) Revoke(tokenHash string) error {
	hashedToken := r.HashToken(tokenHash)
	return r.db.Model(&models.RefreshToken{}).
		Where("token_hash = ?", hashedToken).
		Update("is_revoked", true).Error
}

// RevokeAllForUser revokes all refresh tokens for a user
func (r *RefreshTokenRepository) RevokeAllForUser(userID int64) error {
	return r.db.Model(&models.RefreshToken{}).
		Where("user_id = ?", userID).
		Update("is_revoked", true).Error
}

// RevokeAllForUserExcept revokes all refresh tokens for a user except the specified one
func (r *RefreshTokenRepository) RevokeAllForUserExcept(userID int64, exceptTokenHash string) error {
	hashedToken := r.HashToken(exceptTokenHash)
	return r.db.Model(&models.RefreshToken{}).
		Where("user_id = ? AND token_hash != ?", userID, hashedToken).
		Update("is_revoked", true).Error
}

// DeleteExpired deletes expired refresh tokens
func (r *RefreshTokenRepository) DeleteExpired() error {
	return r.db.Where("expires_at < ?", time.Now().UTC()).
		Delete(&models.RefreshToken{}).Error
}

// SessionRepository handles session data operations
type SessionRepository struct {
	*BaseRepository
}

// NewSessionRepository creates a new SessionRepository
func NewSessionRepository() *SessionRepository {
	return &SessionRepository{
		BaseRepository: NewBaseRepository(),
	}
}

// GenerateSessionID generates a unique session ID
func (r *SessionRepository) GenerateSessionID() string {
	b := make([]byte, 32)
	r.db.Raw("SELECT RANDOMBLOB(32)").Scan(&b)
	return base64.URLEncoding.EncodeToString(b)
}

// Create creates a new session
func (r *SessionRepository) Create(session *models.Session) error {
	return r.db.Create(session).Error
}

// GetBySessionID retrieves a session by session ID
func (r *SessionRepository) GetBySessionID(sessionID string) (*models.Session, error) {
	var session models.Session
	err := r.db.Where("session_id = ? AND is_revoked = false", sessionID).
		Preload("User").
		Preload("Tenant").
		First(&session).Error
	if err != nil {
		return nil, err
	}
	return &session, nil
}

// Revoke revokes a session
func (r *SessionRepository) Revoke(sessionID string) error {
	return r.db.Model(&models.Session{}).
		Where("session_id = ?", sessionID).
		Update("is_revoked", true).Error
}

// RevokeAllForUser revokes all sessions for a user
func (r *SessionRepository) RevokeAllForUser(userID int64) error {
	return r.db.Model(&models.Session{}).
		Where("user_id = ?", userID).
		Update("is_revoked", true).Error
}

// RevokeAllForUserExcept revokes all sessions for a user except the specified one
func (r *SessionRepository) RevokeAllForUserExcept(userID int64, exceptSessionID string) error {
	return r.db.Model(&models.Session{}).
		Where("user_id = ? AND session_id != ?", userID, exceptSessionID).
		Update("is_revoked", true).Error
}

// DeleteExpired deletes expired sessions
func (r *SessionRepository) DeleteExpired() error {
	return r.db.Where("expires_at < ?", time.Now().UTC()).
		Delete(&models.Session{}).Error
}
