package models

import (
	"database/sql/driver"
	"encoding/json"
	"time"
)

// Tenant represents a tenant in the multi-tenant system
type Tenant struct {
	ID          int64     `gorm:"primaryKey;autoIncrement" json:"id"`
	Name        string    `gorm:"size:255;not null" json:"name"`
	Slug        string    `gorm:"size:100;not null;uniqueIndex" json:"slug"`
	Description string    `gorm:"type:text" json:"description"`
	IsActive    bool      `gorm:"default:true" json:"is_active"`
	CreatedAt   time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

// TableName specifies the table name for Tenant
func (Tenant) TableName() string {
	return "tenants"
}

// StringArray is a custom type for storing JSON arrays in MySQL
type StringArray []string

// Scan implements the sql.Scanner interface
func (sa *StringArray) Scan(value interface{}) error {
	if value == nil {
		*sa = []string{}
		return nil
	}
	bytes, ok := value.([]byte)
	if !ok {
		return nil
	}
	return json.Unmarshal(bytes, sa)
}

// Value implements the driver.Valuer interface
func (sa StringArray) Value() (driver.Value, error) {
	if len(sa) == 0 {
		return "[]", nil
	}
	return json.Marshal(sa)
}

// OIDCClient represents an OIDC client application
type OIDCClient struct {
	ID               int64       `gorm:"primaryKey;autoIncrement" json:"id"`
	TenantID         int64       `gorm:"not null;index" json:"tenant_id"`
	ClientID         string      `gorm:"size:255;not null;uniqueIndex" json:"client_id"`
	ClientSecretHash string      `gorm:"size:255;not null" json:"-"`
	Name             string      `gorm:"size:255;not null" json:"name"`
	RedirectURIs     StringArray `gorm:"type:json;not null" json:"redirect_uris"`
	Scopes           StringArray `gorm:"type:json;not null" json:"scopes"`
	GrantTypes       StringArray `gorm:"type:json;not null" json:"grant_types"`
	IsPublic         bool        `gorm:"default:false" json:"is_public"`
	IsActive         bool        `gorm:"default:true" json:"is_active"`
	CreatedAt        time.Time   `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt        time.Time   `gorm:"autoUpdateTime" json:"updated_at"`

	// Relationships
	Tenant *Tenant `gorm:"foreignKey:TenantID" json:"tenant,omitempty"`
}

// TableName specifies the table name for OIDCClient
func (OIDCClient) TableName() string {
	return "oidc_clients"
}

// User represents a user account
type User struct {
	ID                  int64      `gorm:"primaryKey;autoIncrement" json:"id"`
	TenantID            int64      `gorm:"not null;index" json:"tenant_id"`
	Email               string     `gorm:"size:255;not null" json:"email"`
	PasswordHash        string     `gorm:"size:255;not null" json:"-"`
	FullName            string     `gorm:"size:255" json:"full_name"`
	IsActive            bool       `gorm:"default:true;index" json:"is_active"`
	IsEmailVerified     bool       `gorm:"default:false" json:"is_email_verified"`
	FailedLoginAttempts int        `gorm:"default:0" json:"failed_login_attempts"`
	LockedUntil         *time.Time `json:"locked_until"`
	LastLoginAt         *time.Time `json:"last_login_at"`
	CreatedAt           time.Time  `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt           time.Time  `gorm:"autoUpdateTime" json:"updated_at"`

	// Relationships
	Tenant *Tenant `gorm:"foreignKey:TenantID" json:"tenant,omitempty"`
	Roles  []Role  `gorm:"many2many:user_roles" json:"roles,omitempty"`
}

// TableName specifies the table name for User
func (User) TableName() string {
	return "users"
}

// IsLocked checks if the user account is currently locked
func (u *User) IsLocked() bool {
	if u.LockedUntil == nil {
		return false
	}
	return time.Now().Before(*u.LockedUntil)
}

// Role represents an identity role
type Role struct {
	ID           int64     `gorm:"primaryKey;autoIncrement" json:"id"`
	TenantID     int64     `gorm:"not null;index" json:"tenant_id"`
	Name         string    `gorm:"size:100;not null" json:"name"`
	Description  string    `gorm:"type:text" json:"description"`
	IsSystemRole bool      `gorm:"default:false" json:"is_system_role"`
	CreatedAt    time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt    time.Time `gorm:"autoUpdateTime" json:"updated_at"`

	// Relationships
	Tenant *Tenant `gorm:"foreignKey:TenantID" json:"tenant,omitempty"`
	Users  []User  `gorm:"many2many:user_roles" json:"users,omitempty"`
}

// TableName specifies the table name for Role
func (Role) TableName() string {
	return "roles"
}

// UserRole represents the junction table between users and roles
type UserRole struct {
	ID        int64     `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID    int64     `gorm:"not null;index:idx_user_role" json:"user_id"`
	RoleID    int64     `gorm:"not null;index:idx_user_role" json:"role_id"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`

	// Relationships
	User *User `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Role *Role `gorm:"foreignKey:RoleID" json:"role,omitempty"`
}

// TableName specifies the table name for UserRole
func (UserRole) TableName() string {
	return "user_roles"
}

// RefreshToken represents a refresh token for token rotation
type RefreshToken struct {
	ID        int64     `gorm:"primaryKey;autoIncrement" json:"id"`
	TenantID  int64     `gorm:"not null;index" json:"tenant_id"`
	UserID    int64     `gorm:"not null;index" json:"user_id"`
	ClientID  string    `gorm:"size:255;not null" json:"client_id"`
	TokenHash string    `gorm:"size:255;not null;uniqueIndex" json:"-"`
	IsRevoked bool      `gorm:"default:false;index" json:"is_revoked"`
	ExpiresAt time.Time `gorm:"not null;index" json:"expires_at"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`

	// Relationships
	Tenant *Tenant `gorm:"foreignKey:TenantID" json:"tenant,omitempty"`
	User   *User   `gorm:"foreignKey:UserID" json:"user,omitempty"`
}

// TableName specifies the table name for RefreshToken
func (RefreshToken) TableName() string {
	return "refresh_tokens"
}

// IsExpired checks if the refresh token is expired
func (rt *RefreshToken) IsExpired() bool {
	return time.Now().After(rt.ExpiresAt)
}

// AuthorizationCode represents an authorization code for OIDC flow
type AuthorizationCode struct {
	ID                  int64     `gorm:"primaryKey;autoIncrement" json:"id"`
	TenantID            int64     `gorm:"not null" json:"tenant_id"`
	UserID              int64     `gorm:"not null" json:"user_id"`
	ClientID            string    `gorm:"size:255;not null" json:"client_id"`
	Code                string    `gorm:"size:255;not null;uniqueIndex" json:"-"`
	RedirectURI         string    `gorm:"size:500;not null" json:"redirect_uri"`
	Scope               string    `gorm:"size:500" json:"scope"`
	CodeChallenge       string    `gorm:"size:255" json:"code_challenge"`
	CodeChallengeMethod string    `gorm:"size:10" json:"code_challenge_method"`
	ExpiresAt           time.Time `gorm:"not null;index" json:"expires_at"`
	IsUsed              bool      `gorm:"default:false;index" json:"is_used"`
	CreatedAt           time.Time `gorm:"autoCreateTime" json:"created_at"`

	// Relationships
	Tenant *Tenant `gorm:"foreignKey:TenantID" json:"tenant,omitempty"`
	User   *User   `gorm:"foreignKey:UserID" json:"user,omitempty"`
}

// TableName specifies the table name for AuthorizationCode
func (AuthorizationCode) TableName() string {
	return "authorization_codes"
}

// IsExpired checks if the authorization code is expired
func (ac *AuthorizationCode) IsExpired() bool {
	return time.Now().After(ac.ExpiresAt)
}

// Session represents a user session for single sign-out
type Session struct {
	ID        int64     `gorm:"primaryKey;autoIncrement" json:"id"`
	TenantID  int64     `gorm:"not null;index" json:"tenant_id"`
	UserID    int64     `gorm:"not null;index" json:"user_id"`
	ClientID  string    `gorm:"size:255;not null" json:"client_id"`
	SessionID string    `gorm:"size:255;not null;uniqueIndex" json:"session_id"`
	IsRevoked bool      `gorm:"default:false;index" json:"is_revoked"`
	ExpiresAt time.Time `gorm:"not null;index" json:"expires_at"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`

	// Relationships
	Tenant *Tenant `gorm:"foreignKey:TenantID" json:"tenant,omitempty"`
	User   *User   `gorm:"foreignKey:UserID" json:"user,omitempty"`
}

// TableName specifies the table name for Session
func (Session) TableName() string {
	return "sessions"
}

// IsExpired checks if the session is expired
func (s *Session) IsExpired() bool {
	return time.Now().After(s.ExpiresAt)
}
