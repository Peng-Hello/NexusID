package dto

// TenantDTO represents a tenant data transfer object
type TenantDTO struct {
	ID          int64  `json:"id"`
	Name        string `json:"name"`
	Slug        string `json:"slug"`
	Description string `json:"description"`
	IsActive    bool   `json:"is_active"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
}

// CreateTenantRequest represents a request to create a tenant
type CreateTenantRequest struct {
	Name        string `json:"name" binding:"required"`
	Slug        string `json:"slug" binding:"required"`
	Description string `json:"description"`
}

// UpdateTenantRequest represents a request to update a tenant
type UpdateTenantRequest struct {
	Name        *string `json:"name"`
	Slug        *string `json:"slug"`
	Description *string `json:"description"`
	IsActive    *bool   `json:"is_active"`
}

// OIDCClientDTO represents an OIDC client data transfer object
type OIDCClientDTO struct {
	ID           int64    `json:"id"`
	TenantID     int64    `json:"tenant_id"`
	ClientID     string   `json:"client_id"`
	Name         string   `json:"name"`
	RedirectURIs []string `json:"redirect_uris"`
	Scopes       []string `json:"scopes"`
	GrantTypes   []string `json:"grant_types"`
	IsPublic     bool     `json:"is_public"`
	IsActive     bool     `json:"is_active"`
	CreatedAt    string   `json:"created_at"`
	UpdatedAt    string   `json:"updated_at"`
}

// CreateOIDCClientRequest represents a request to create an OIDC client
type CreateOIDCClientRequest struct {
	Name         string   `json:"name" binding:"required"`
	RedirectURIs []string `json:"redirect_uris" binding:"required,min=1"`
	Scopes       []string `json:"scopes" binding:"required,min=1"`
	GrantTypes   []string `json:"grant_types" binding:"required,min=1"`
	IsPublic     bool     `json:"is_public"`
}

// UpdateOIDCClientRequest represents a request to update an OIDC client
type UpdateOIDCClientRequest struct {
	Name         *string   `json:"name"`
	RedirectURIs *[]string `json:"redirect_uris"`
	Scopes       *[]string `json:"scopes"`
	GrantTypes   *[]string `json:"grant_types"`
	IsPublic     *bool     `json:"is_public"`
	IsActive     *bool     `json:"is_active"`
}

// UserDTO represents a user data transfer object
type UserDTO struct {
	ID              int64    `json:"id"`
	TenantID        int64    `json:"tenant_id"`
	Email           string   `json:"email"`
	FullName        string   `json:"full_name"`
	IsActive        bool     `json:"is_active"`
	IsEmailVerified bool     `json:"is_email_verified"`
	Roles           []string `json:"roles,omitempty"`
	LockedUntil     *string  `json:"locked_until,omitempty"`
	LastLoginAt     *string  `json:"last_login_at,omitempty"`
	CreatedAt       string   `json:"created_at"`
	UpdatedAt       string   `json:"updated_at"`
}

// CreateUserRequest represents a request to create a user
type CreateUserRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
	FullName string `json:"full_name"`
}

// UpdateUserRequest represents a request to update a user
type UpdateUserRequest struct {
	FullName        *string `json:"full_name"`
	IsActive        *bool   `json:"is_active"`
	IsEmailVerified *bool   `json:"is_email_verified"`
}

// LoginRequest represents a login request
type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

// LoginResponse represents a login response
type LoginResponse struct {
	User         UserDTO `json:"user"`
	AccessToken  string  `json:"access_token"`
	RefreshToken string  `json:"refresh_token"`
	ExpiresIn    int64   `json:"expires_in"`
}

// RoleDTO represents a role data transfer object
type RoleDTO struct {
	ID           int64  `json:"id"`
	TenantID     int64  `json:"tenant_id"`
	Name         string `json:"name"`
	Description  string `json:"description"`
	IsSystemRole bool   `json:"is_system_role"`
	CreatedAt    string `json:"created_at"`
	UpdatedAt    string `json:"updated_at"`
}

// CreateRoleRequest represents a request to create a role
type CreateRoleRequest struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
}

// UpdateRoleRequest represents a request to update a role
type UpdateRoleRequest struct {
	Name        *string `json:"name"`
	Description *string `json:"description"`
}

// AssignRoleRequest represents a request to assign a role to a user
type AssignRoleRequest struct {
	UserID int64 `json:"user_id" binding:"required"`
	RoleID int64 `json:"role_id" binding:"required"`
}

// ChangePasswordRequest represents a request to change password
type ChangePasswordRequest struct {
	NewPassword string `json:"new_password" binding:"required,min=8"`
}
