package service

import (
	"errors"
	"fmt"

	"github.com/nexus-id/backend/internal/dto"
	"github.com/nexus-id/backend/internal/models"
	"github.com/nexus-id/backend/internal/repository"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

var (
	ErrTenantNotFound     = errors.New("tenant not found")
	ErrTenantSlugExists   = errors.New("tenant slug already exists")
	ErrUserNotFound       = errors.New("user not found")
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrAccountLocked      = errors.New("account is locked")
)

// TenantService handles tenant business logic
type TenantService struct {
	tenantRepo *repository.TenantRepository
}

// NewTenantService creates a new TenantService
func NewTenantService() *TenantService {
	return &TenantService{
		tenantRepo: repository.NewTenantRepository(),
	}
}

// Create creates a new tenant
func (s *TenantService) Create(req *dto.CreateTenantRequest) (*dto.TenantDTO, error) {
	// Check if slug already exists
	if _, err := s.tenantRepo.GetBySlug(req.Slug); err == nil {
		return nil, ErrTenantSlugExists
	}

	tenant := &models.Tenant{
		Name:        req.Name,
		Slug:        req.Slug,
		Description: req.Description,
		IsActive:    true,
	}

	if err := s.tenantRepo.Create(tenant); err != nil {
		return nil, fmt.Errorf("failed to create tenant: %w", err)
	}

	return s.toDTO(tenant), nil
}

// GetByID retrieves a tenant by ID
func (s *TenantService) GetByID(id int64) (*dto.TenantDTO, error) {
	tenant, err := s.tenantRepo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrTenantNotFound
		}
		return nil, fmt.Errorf("failed to get tenant: %w", err)
	}

	return s.toDTO(tenant), nil
}

// GetBySlug retrieves a tenant by slug
func (s *TenantService) GetBySlug(slug string) (*dto.TenantDTO, error) {
	tenant, err := s.tenantRepo.GetBySlug(slug)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrTenantNotFound
		}
		return nil, fmt.Errorf("failed to get tenant: %w", err)
	}

	return s.toDTO(tenant), nil
}

// List retrieves tenants with pagination
func (s *TenantService) List(page, pageSize int) ([]*dto.TenantDTO, int64, error) {
	offset := (page - 1) * pageSize
	tenants, total, err := s.tenantRepo.List(offset, pageSize)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list tenants: %w", err)
	}

	dtos := make([]*dto.TenantDTO, len(tenants))
	for i, tenant := range tenants {
		dtos[i] = s.toDTO(tenant)
	}

	return dtos, total, nil
}

// Update updates a tenant
func (s *TenantService) Update(id int64, req *dto.UpdateTenantRequest) (*dto.TenantDTO, error) {
	tenant, err := s.tenantRepo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrTenantNotFound
		}
		return nil, fmt.Errorf("failed to get tenant: %w", err)
	}

	// Check if new slug already exists (if being updated)
	if req.Slug != nil && *req.Slug != tenant.Slug {
		if _, err := s.tenantRepo.GetBySlug(*req.Slug); err == nil {
			return nil, ErrTenantSlugExists
		}
		tenant.Slug = *req.Slug
	}

	if req.Name != nil {
		tenant.Name = *req.Name
	}
	if req.Description != nil {
		tenant.Description = *req.Description
	}
	if req.IsActive != nil {
		tenant.IsActive = *req.IsActive
	}

	if err := s.tenantRepo.Update(tenant); err != nil {
		return nil, fmt.Errorf("failed to update tenant: %w", err)
	}

	return s.toDTO(tenant), nil
}

// Delete soft deletes a tenant
func (s *TenantService) Delete(id int64) error {
	if err := s.tenantRepo.Delete(id); err != nil {
		return fmt.Errorf("failed to delete tenant: %w", err)
	}
	return nil
}

// toDTO converts a model to DTO
func (s *TenantService) toDTO(tenant *models.Tenant) *dto.TenantDTO {
	return &dto.TenantDTO{
		ID:          tenant.ID,
		Name:        tenant.Name,
		Slug:        tenant.Slug,
		Description: tenant.Description,
		IsActive:    tenant.IsActive,
		CreatedAt:   tenant.CreatedAt.Format("2006-01-02T15:04:05Z"),
		UpdatedAt:   tenant.UpdatedAt.Format("2006-01-02T15:04:05Z"),
	}
}

// UserService handles user business logic
type UserService struct {
	userRepo   *repository.UserRepository
	roleRepo   *repository.RoleRepository
	tenantRepo *repository.TenantRepository
}

// NewUserService creates a new UserService
func NewUserService() *UserService {
	return &UserService{
		userRepo:   repository.NewUserRepository(),
		roleRepo:   repository.NewRoleRepository(),
		tenantRepo: repository.NewTenantRepository(),
	}
}

// Create creates a new user
func (s *UserService) Create(tenantID int64, req *dto.CreateUserRequest) (*dto.UserDTO, error) {
	// Verify tenant exists
	if _, err := s.tenantRepo.GetByID(tenantID); err != nil {
		return nil, ErrTenantNotFound
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	user := &models.User{
		TenantID:     tenantID,
		Email:        req.Email,
		PasswordHash: string(hashedPassword),
		FullName:     req.FullName,
		IsActive:     true,
	}

	if err := s.userRepo.Create(user); err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	return s.toDTO(user), nil
}

// GetByID retrieves a user by ID
func (s *UserService) GetByID(tenantID, userID int64) (*dto.UserDTO, error) {
	user, err := s.userRepo.GetByID(userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrUserNotFound
		}
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	// Ensure user belongs to tenant
	if user.TenantID != tenantID {
		return nil, ErrUserNotFound
	}

	return s.toDTO(user), nil
}

// GetByIDModel retrieves a user model by ID (internal use)
func (s *UserService) GetByIDModel(tenantID, userID int64) (*models.User, error) {
	user, err := s.userRepo.GetByID(userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrUserNotFound
		}
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	// Ensure user belongs to tenant
	if user.TenantID != tenantID {
		return nil, ErrUserNotFound
	}

	return user, nil
}

// GetByEmail retrieves a user by email within a tenant
func (s *UserService) GetByEmail(tenantID int64, email string) (*models.User, error) {
	user, err := s.userRepo.GetByEmail(tenantID, email)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrUserNotFound
		}
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	return user, nil
}

// List retrieves users for a tenant with pagination
func (s *UserService) List(tenantID int64, page, pageSize int) ([]*dto.UserDTO, int64, error) {
	offset := (page - 1) * pageSize
	users, total, err := s.userRepo.List(tenantID, offset, pageSize)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list users: %w", err)
	}

	dtos := make([]*dto.UserDTO, len(users))
	for i, user := range users {
		dtos[i] = s.toDTO(user)
	}

	return dtos, total, nil
}

// Update updates a user
func (s *UserService) Update(tenantID, userID int64, req *dto.UpdateUserRequest) (*dto.UserDTO, error) {
	user, err := s.userRepo.GetByID(userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrUserNotFound
		}
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	// Ensure user belongs to tenant
	if user.TenantID != tenantID {
		return nil, ErrUserNotFound
	}

	if req.FullName != nil {
		user.FullName = *req.FullName
	}
	if req.IsActive != nil {
		user.IsActive = *req.IsActive
	}
	if req.IsEmailVerified != nil {
		user.IsEmailVerified = *req.IsEmailVerified
	}

	if err := s.userRepo.Update(user); err != nil {
		return nil, fmt.Errorf("failed to update user: %w", err)
	}

	return s.toDTO(user), nil
}

// Delete soft deletes a user
func (s *UserService) Delete(tenantID, userID int64) error {
	user, err := s.userRepo.GetByID(userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrUserNotFound
		}
		return fmt.Errorf("failed to get user: %w", err)
	}

	// Ensure user belongs to tenant
	if user.TenantID != tenantID {
		return ErrUserNotFound
	}

	if err := s.userRepo.Delete(userID); err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}

	return nil
}

// ValidatePassword validates a user's password
func (s *UserService) ValidatePassword(tenantID int64, email, password string) (*models.User, error) {
	user, err := s.GetByEmail(tenantID, email)
	if err != nil {
		return nil, ErrInvalidCredentials
	}

	// Check if account is locked
	if user.IsLocked() {
		return nil, ErrAccountLocked
	}

	// Check if account is active
	if !user.IsActive {
		return nil, ErrInvalidCredentials
	}

	// Verify password
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		// Increment failed login attempts
		s.userRepo.IncrementFailedLoginAttempts(user.ID)
		return nil, ErrInvalidCredentials
	}

	// Reset failed login attempts on successful login
	s.userRepo.ResetFailedLoginAttempts(user.ID)

	return user, nil
}

// UpdateLastLogin updates the user's last login time
func (s *UserService) UpdateLastLogin(userID int64) error {
	return s.userRepo.UpdateLastLogin(userID)
}

// toDTO converts a user model to DTO
func (s *UserService) toDTO(user *models.User) *dto.UserDTO {
	roles := make([]string, len(user.Roles))
	for i, role := range user.Roles {
		roles[i] = role.Name
	}

	var lockedUntil *string
	if user.LockedUntil != nil {
		s := user.LockedUntil.Format("2006-01-02T15:04:05Z")
		lockedUntil = &s
	}

	var lastLoginAt *string
	if user.LastLoginAt != nil {
		s := user.LastLoginAt.Format("2006-01-02T15:04:05Z")
		lastLoginAt = &s
	}

	return &dto.UserDTO{
		ID:              user.ID,
		TenantID:        user.TenantID,
		Email:           user.Email,
		FullName:        user.FullName,
		IsActive:        user.IsActive,
		IsEmailVerified: user.IsEmailVerified,
		Roles:           roles,
		LockedUntil:     lockedUntil,
		LastLoginAt:     lastLoginAt,
		CreatedAt:       user.CreatedAt.Format("2006-01-02T15:04:05Z"),
		UpdatedAt:       user.UpdatedAt.Format("2006-01-02T15:04:05Z"),
	}
}
