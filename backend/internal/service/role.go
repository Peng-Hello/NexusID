package service

import (
	"errors"
	"fmt"

	"github.com/nexus-id/backend/internal/dto"
	"github.com/nexus-id/backend/internal/models"
	"github.com/nexus-id/backend/internal/repository"
	"gorm.io/gorm"
)

// RoleService handles role business logic
type RoleService struct {
	roleRepo   *repository.RoleRepository
	tenantRepo *repository.TenantRepository
}

// NewRoleService creates a new RoleService
func NewRoleService() *RoleService {
	return &RoleService{
		roleRepo:   repository.NewRoleRepository(),
		tenantRepo: repository.NewTenantRepository(),
	}
}

// Create creates a new role
func (s *RoleService) Create(tenantID int64, req *dto.CreateRoleRequest) (*dto.RoleDTO, error) {
	// Verify tenant exists
	if _, err := s.tenantRepo.GetByID(tenantID); err != nil {
		return nil, ErrTenantNotFound
	}

	// Check if role name already exists for this tenant
	if _, err := s.roleRepo.GetByName(tenantID, req.Name); err == nil {
		return nil, errors.New("role name already exists")
	}

	role := &models.Role{
		TenantID:     tenantID,
		Name:         req.Name,
		Description:  req.Description,
		IsSystemRole: false,
	}

	if err := s.roleRepo.Create(role); err != nil {
		return nil, fmt.Errorf("failed to create role: %w", err)
	}

	return s.toDTO(role), nil
}

// GetByID retrieves a role by ID
func (s *RoleService) GetByID(tenantID, roleID int64) (*dto.RoleDTO, error) {
	role, err := s.roleRepo.GetByID(roleID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("role not found")
		}
		return nil, fmt.Errorf("failed to get role: %w", err)
	}

	// Ensure role belongs to tenant
	if role.TenantID != tenantID {
		return nil, errors.New("role not found")
	}

	return s.toDTO(role), nil
}

// List retrieves roles for a tenant with pagination
func (s *RoleService) List(tenantID int64, page, pageSize int) ([]*dto.RoleDTO, int64, error) {
	offset := (page - 1) * pageSize
	roles, total, err := s.roleRepo.List(tenantID, offset, pageSize)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list roles: %w", err)
	}

	dtos := make([]*dto.RoleDTO, len(roles))
	for i, role := range roles {
		dtos[i] = s.toDTO(role)
	}

	return dtos, total, nil
}

// Update updates a role
func (s *RoleService) Update(tenantID, roleID int64, req *dto.UpdateRoleRequest) (*dto.RoleDTO, error) {
	role, err := s.roleRepo.GetByID(roleID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("role not found")
		}
		return nil, fmt.Errorf("failed to get role: %w", err)
	}

	// Ensure role belongs to tenant
	if role.TenantID != tenantID {
		return nil, errors.New("role not found")
	}

	// Cannot modify system roles
	if role.IsSystemRole {
		return nil, errors.New("cannot modify system roles")
	}

	// Check if new name already exists (if being updated)
	if req.Name != nil && *req.Name != role.Name {
		if _, err := s.roleRepo.GetByName(tenantID, *req.Name); err == nil {
			return nil, errors.New("role name already exists")
		}
		role.Name = *req.Name
	}

	if req.Description != nil {
		role.Description = *req.Description
	}

	if err := s.roleRepo.Update(role); err != nil {
		return nil, fmt.Errorf("failed to update role: %w", err)
	}

	return s.toDTO(role), nil
}

// Delete deletes a role
func (s *RoleService) Delete(tenantID, roleID int64) error {
	role, err := s.roleRepo.GetByID(roleID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("role not found")
		}
		return fmt.Errorf("failed to get role: %w", err)
	}

	// Ensure role belongs to tenant
	if role.TenantID != tenantID {
		return errors.New("role not found")
	}

	// Cannot delete system roles
	if role.IsSystemRole {
		return errors.New("cannot delete system roles")
	}

	if err := s.roleRepo.Delete(roleID); err != nil {
		return fmt.Errorf("failed to delete role: %w", err)
	}

	return nil
}

// AssignRole assigns a role to a user
func (s *RoleService) AssignRole(tenantID, userID, roleID int64) error {
	// Verify role exists and belongs to tenant
	role, err := s.roleRepo.GetByID(roleID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("role not found")
		}
		return fmt.Errorf("failed to get role: %w", err)
	}

	if role.TenantID != tenantID {
		return errors.New("role not found")
	}

	// Assign role
	if err := s.roleRepo.AssignRole(userID, roleID); err != nil {
		return fmt.Errorf("failed to assign role: %w", err)
	}

	return nil
}

// RevokeRole revokes a role from a user
func (s *RoleService) RevokeRole(tenantID, userID, roleID int64) error {
	// Verify role exists and belongs to tenant
	role, err := s.roleRepo.GetByID(roleID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("role not found")
		}
		return fmt.Errorf("failed to get role: %w", err)
	}

	if role.TenantID != tenantID {
		return errors.New("role not found")
	}

	// Revoke role
	if err := s.roleRepo.RevokeRole(userID, roleID); err != nil {
		return fmt.Errorf("failed to revoke role: %w", err)
	}

	return nil
}

// GetUserRoles retrieves all roles for a user
func (s *RoleService) GetUserRoles(userID int64) ([]*dto.RoleDTO, error) {
	roles, err := s.roleRepo.GetUserRoles(userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user roles: %w", err)
	}

	dtos := make([]*dto.RoleDTO, len(roles))
	for i, role := range roles {
		dtos[i] = s.toDTO(role)
	}

	return dtos, nil
}

// toDTO converts a model to DTO
func (s *RoleService) toDTO(role *models.Role) *dto.RoleDTO {
	return &dto.RoleDTO{
		ID:           role.ID,
		TenantID:     role.TenantID,
		Name:         role.Name,
		Description:  role.Description,
		IsSystemRole: role.IsSystemRole,
		CreatedAt:    role.CreatedAt.Format("2006-01-02T15:04:05Z"),
		UpdatedAt:    role.UpdatedAt.Format("2006-01-02T15:04:05Z"),
	}
}
