package repository

import (
	"github.com/nexus-id/backend/internal/database"
	"github.com/nexus-id/backend/internal/models"
	"gorm.io/gorm"
)

// BaseRepository provides common database operations
type BaseRepository struct {
	db *gorm.DB
}

// NewBaseRepository creates a new base repository
func NewBaseRepository() *BaseRepository {
	return &BaseRepository{
		db: database.DB,
	}
}

// GetDB returns the database connection
func (r *BaseRepository) GetDB() *gorm.DB {
	return r.db
}

// TenantRepository handles tenant data operations
type TenantRepository struct {
	*BaseRepository
}

// NewTenantRepository creates a new tenant repository
func NewTenantRepository() *TenantRepository {
	return &TenantRepository{
		BaseRepository: NewBaseRepository(),
	}
}

// Create creates a new tenant
func (r *TenantRepository) Create(tenant *models.Tenant) error {
	return r.db.Create(tenant).Error
}

// GetByID retrieves a tenant by ID
func (r *TenantRepository) GetByID(id int64) (*models.Tenant, error) {
	var tenant models.Tenant
	err := r.db.Where("id = ?", id).First(&tenant).Error
	if err != nil {
		return nil, err
	}
	return &tenant, nil
}

// GetBySlug retrieves a tenant by slug
func (r *TenantRepository) GetBySlug(slug string) (*models.Tenant, error) {
	var tenant models.Tenant
	err := r.db.Where("slug = ?", slug).First(&tenant).Error
	if err != nil {
		return nil, err
	}
	return &tenant, nil
}

// List retrieves all tenants with pagination
func (r *TenantRepository) List(offset, limit int) ([]*models.Tenant, int64, error) {
	var tenants []*models.Tenant
	var total int64

	if err := r.db.Model(&models.Tenant{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := r.db.Offset(offset).Limit(limit).Find(&tenants).Error
	return tenants, total, err
}

// Update updates a tenant
func (r *TenantRepository) Update(tenant *models.Tenant) error {
	return r.db.Save(tenant).Error
}

// Delete soft deletes a tenant
func (r *TenantRepository) Delete(id int64) error {
	return r.db.Delete(&models.Tenant{}, id).Error
}

// UserRepository handles user data operations
type UserRepository struct {
	*BaseRepository
}

// NewUserRepository creates a new user repository
func NewUserRepository() *UserRepository {
	return &UserRepository{
		BaseRepository: NewBaseRepository(),
	}
}

// Create creates a new user
func (r *UserRepository) Create(user *models.User) error {
	return r.db.Create(user).Error
}

// GetByID retrieves a user by ID
func (r *UserRepository) GetByID(id int64) (*models.User, error) {
	var user models.User
	err := r.db.Preload("Roles").Where("id = ?", id).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// GetByEmail retrieves a user by email within a tenant
func (r *UserRepository) GetByEmail(tenantID int64, email string) (*models.User, error) {
	var user models.User
	err := r.db.Preload("Roles").Where("tenant_id = ? AND email = ?", tenantID, email).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// List retrieves users for a tenant with pagination
func (r *UserRepository) List(tenantID int64, offset, limit int) ([]*models.User, int64, error) {
	var users []*models.User
	var total int64

	query := r.db.Model(&models.User{}).Where("tenant_id = ?", tenantID)
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := query.Preload("Roles").Offset(offset).Limit(limit).Find(&users).Error
	return users, total, err
}

// Update updates a user
func (r *UserRepository) Update(user *models.User) error {
	return r.db.Save(user).Error
}

// Delete soft deletes a user
func (r *UserRepository) Delete(id int64) error {
	return r.db.Delete(&models.User{}, id).Error
}

// UpdateLastLogin updates the user's last login time
func (r *UserRepository) UpdateLastLogin(userID int64) error {
	return r.db.Model(&models.User{}).Where("id = ?", userID).Update("last_login_at", gorm.Expr("NOW()")).Error
}

// IncrementFailedLoginAttempts increments the failed login attempts counter
func (r *UserRepository) IncrementFailedLoginAttempts(userID int64) error {
	return r.db.Model(&models.User{}).Where("id = ?", userID).
		UpdateColumn("failed_login_attempts", gorm.Expr("failed_login_attempts + 1")).Error
}

// ResetFailedLoginAttempts resets the failed login attempts counter
func (r *UserRepository) ResetFailedLoginAttempts(userID int64) error {
	return r.db.Model(&models.User{}).Where("id = ?", userID).
		Update("failed_login_attempts", 0).Error
}

// LockAccount locks a user account until a specified time
func (r *UserRepository) LockAccount(userID int64, lockedUntil interface{}) error {
	return r.db.Model(&models.User{}).Where("id = ?", userID).
		Update("locked_until", lockedUntil).Error
}

// OIDCClientRepository handles OIDC client data operations
type OIDCClientRepository struct {
	*BaseRepository
}

// NewOIDCClientRepository creates a new OIDC client repository
func NewOIDCClientRepository() *OIDCClientRepository {
	return &OIDCClientRepository{
		BaseRepository: NewBaseRepository(),
	}
}

// Create creates a new OIDC client
func (r *OIDCClientRepository) Create(client *models.OIDCClient) error {
	return r.db.Create(client).Error
}

// GetByID retrieves an OIDC client by ID
func (r *OIDCClientRepository) GetByID(id int64) (*models.OIDCClient, error) {
	var client models.OIDCClient
	err := r.db.Where("id = ?", id).First(&client).Error
	if err != nil {
		return nil, err
	}
	return &client, nil
}

// GetByClientID retrieves an OIDC client by client_id
func (r *OIDCClientRepository) GetByClientID(clientID string) (*models.OIDCClient, error) {
	var client models.OIDCClient
	err := r.db.Where("client_id = ?", clientID).First(&client).Error
	if err != nil {
		return nil, err
	}
	return &client, nil
}

// List retrieves OIDC clients for a tenant with pagination
func (r *OIDCClientRepository) List(tenantID int64, offset, limit int) ([]*models.OIDCClient, int64, error) {
	var clients []*models.OIDCClient
	var total int64

	query := r.db.Model(&models.OIDCClient{}).Where("tenant_id = ?", tenantID)
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := query.Offset(offset).Limit(limit).Find(&clients).Error
	return clients, total, err
}

// Update updates an OIDC client
func (r *OIDCClientRepository) Update(client *models.OIDCClient) error {
	return r.db.Save(client).Error
}

// Delete soft deletes an OIDC client
func (r *OIDCClientRepository) Delete(id int64) error {
	return r.db.Delete(&models.OIDCClient{}, id).Error
}

// RoleRepository handles role data operations
type RoleRepository struct {
	*BaseRepository
}

// NewRoleRepository creates a new role repository
func NewRoleRepository() *RoleRepository {
	return &RoleRepository{
		BaseRepository: NewBaseRepository(),
	}
}

// Create creates a new role
func (r *RoleRepository) Create(role *models.Role) error {
	return r.db.Create(role).Error
}

// GetByID retrieves a role by ID
func (r *RoleRepository) GetByID(id int64) (*models.Role, error) {
	var role models.Role
	err := r.db.Where("id = ?", id).First(&role).Error
	if err != nil {
		return nil, err
	}
	return &role, nil
}

// GetByName retrieves a role by name within a tenant
func (r *RoleRepository) GetByName(tenantID int64, name string) (*models.Role, error) {
	var role models.Role
	err := r.db.Where("tenant_id = ? AND name = ?", tenantID, name).First(&role).Error
	if err != nil {
		return nil, err
	}
	return &role, nil
}

// List retrieves roles for a tenant with pagination
func (r *RoleRepository) List(tenantID int64, offset, limit int) ([]*models.Role, int64, error) {
	var roles []*models.Role
	var total int64

	query := r.db.Model(&models.Role{}).Where("tenant_id = ?", tenantID)
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := query.Offset(offset).Limit(limit).Find(&roles).Error
	return roles, total, err
}

// Update updates a role
func (r *RoleRepository) Update(role *models.Role) error {
	return r.db.Save(role).Error
}

// Delete soft deletes a role
func (r *RoleRepository) Delete(id int64) error {
	return r.db.Delete(&models.Role{}, id).Error
}

// AssignRole assigns a role to a user
func (r *RoleRepository) AssignRole(userID, roleID int64) error {
	userRole := &models.UserRole{
		UserID: userID,
		RoleID: roleID,
	}
	return r.db.Create(userRole).Error
}

// RevokeRole revokes a role from a user
func (r *RoleRepository) RevokeRole(userID, roleID int64) error {
	return r.db.Where("user_id = ? AND role_id = ?", userID, roleID).Delete(&models.UserRole{}).Error
}

// GetUserRoles retrieves all roles for a user
func (r *RoleRepository) GetUserRoles(userID int64) ([]*models.Role, error) {
	var roles []*models.Role
	err := r.db.Joins("JOIN user_roles ON user_roles.role_id = roles.id").
		Where("user_roles.user_id = ?", userID).
		Find(&roles).Error
	return roles, err
}
