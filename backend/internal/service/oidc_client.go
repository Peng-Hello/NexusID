package service

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"

	"github.com/nexus-id/backend/internal/dto"
	"github.com/nexus-id/backend/internal/models"
	"github.com/nexus-id/backend/internal/repository"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

var (
	ErrClientNotFound = errors.New("client not found")
	ErrClientIDExists = errors.New("client ID already exists")
)

// OIDCClientService handles OIDC client business logic
type OIDCClientService struct {
	clientRepo *repository.OIDCClientRepository
	tenantRepo *repository.TenantRepository
}

// NewOIDCClientService creates a new OIDCClientService
func NewOIDCClientService() *OIDCClientService {
	return &OIDCClientService{
		clientRepo: repository.NewOIDCClientRepository(),
		tenantRepo: repository.NewTenantRepository(),
	}
}

// generateClientID generates a unique client ID
func generateClientID() (string, error) {
	b := make([]byte, 16)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return "nexus_" + base64.URLEncoding.EncodeToString(b), nil
}

// generateClientSecret generates a random client secret
func generateClientSecret() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(b), nil
}

// Create creates a new OIDC client
func (s *OIDCClientService) Create(tenantID int64, req *dto.CreateOIDCClientRequest) (*dto.OIDCClientDTO, string, error) {
	// Verify tenant exists
	if _, err := s.tenantRepo.GetByID(tenantID); err != nil {
		return nil, "", ErrTenantNotFound
	}

	// Generate client ID and secret
	clientID, err := generateClientID()
	if err != nil {
		return nil, "", fmt.Errorf("failed to generate client ID: %w", err)
	}

	clientSecret, err := generateClientSecret()
	if err != nil {
		return nil, "", fmt.Errorf("failed to generate client secret: %w", err)
	}

	// Hash client secret for storage
	hashedSecret, err := bcrypt.GenerateFromPassword([]byte(clientSecret), bcrypt.DefaultCost)
	if err != nil {
		return nil, "", fmt.Errorf("failed to hash client secret: %w", err)
	}

	client := &models.OIDCClient{
		TenantID:         tenantID,
		ClientID:         clientID,
		ClientSecretHash: string(hashedSecret),
		Name:             req.Name,
		RedirectURIs:     req.RedirectURIs,
		Scopes:           req.Scopes,
		GrantTypes:       req.GrantTypes,
		IsPublic:         req.IsPublic,
		IsActive:         true,
	}

	if err := s.clientRepo.Create(client); err != nil {
		return nil, "", fmt.Errorf("failed to create client: %w", err)
	}

	return s.toDTO(client), clientSecret, nil
}

// GetByID retrieves an OIDC client by ID
func (s *OIDCClientService) GetByID(tenantID, clientID int64) (*dto.OIDCClientDTO, error) {
	client, err := s.clientRepo.GetByID(clientID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrClientNotFound
		}
		return nil, fmt.Errorf("failed to get client: %w", err)
	}

	// Ensure client belongs to tenant
	if client.TenantID != tenantID {
		return nil, ErrClientNotFound
	}

	return s.toDTO(client), nil
}

// GetByClientID retrieves an OIDC client by client_id
func (s *OIDCClientService) GetByClientID(clientID string) (*models.OIDCClient, error) {
	client, err := s.clientRepo.GetByClientID(clientID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrClientNotFound
		}
		return nil, fmt.Errorf("failed to get client: %w", err)
	}

	return client, nil
}

// ValidateClientSecret validates a client's secret
func (s *OIDCClientService) ValidateClientSecret(clientID, clientSecret string) (*models.OIDCClient, error) {
	client, err := s.GetByClientID(clientID)
	if err != nil {
		return nil, ErrInvalidCredentials
	}

	// Check if client is active
	if !client.IsActive {
		return nil, ErrInvalidCredentials
	}

	// Verify secret (skip for public clients)
	if !client.IsPublic {
		if err := bcrypt.CompareHashAndPassword([]byte(client.ClientSecretHash), []byte(clientSecret)); err != nil {
			return nil, ErrInvalidCredentials
		}
	}

	return client, nil
}

// List retrieves OIDC clients for a tenant with pagination
func (s *OIDCClientService) List(tenantID int64, page, pageSize int) ([]*dto.OIDCClientDTO, int64, error) {
	offset := (page - 1) * pageSize
	clients, total, err := s.clientRepo.List(tenantID, offset, pageSize)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list clients: %w", err)
	}

	dtos := make([]*dto.OIDCClientDTO, len(clients))
	for i, client := range clients {
		dtos[i] = s.toDTO(client)
	}

	return dtos, total, nil
}

// Update updates an OIDC client
func (s *OIDCClientService) Update(tenantID, clientID int64, req *dto.UpdateOIDCClientRequest) (*dto.OIDCClientDTO, error) {
	client, err := s.clientRepo.GetByID(clientID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrClientNotFound
		}
		return nil, fmt.Errorf("failed to get client: %w", err)
	}

	// Ensure client belongs to tenant
	if client.TenantID != tenantID {
		return nil, ErrClientNotFound
	}

	if req.Name != nil {
		client.Name = *req.Name
	}
	if req.RedirectURIs != nil {
		client.RedirectURIs = *req.RedirectURIs
	}
	if req.Scopes != nil {
		client.Scopes = *req.Scopes
	}
	if req.GrantTypes != nil {
		client.GrantTypes = *req.GrantTypes
	}
	if req.IsPublic != nil {
		client.IsPublic = *req.IsPublic
	}
	if req.IsActive != nil {
		client.IsActive = *req.IsActive
	}

	if err := s.clientRepo.Update(client); err != nil {
		return nil, fmt.Errorf("failed to update client: %w", err)
	}

	return s.toDTO(client), nil
}

// Delete soft deletes an OIDC client
func (s *OIDCClientService) Delete(tenantID, clientID int64) error {
	client, err := s.clientRepo.GetByID(clientID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrClientNotFound
		}
		return fmt.Errorf("failed to get client: %w", err)
	}

	// Ensure client belongs to tenant
	if client.TenantID != tenantID {
		return ErrClientNotFound
	}

	if err := s.clientRepo.Delete(clientID); err != nil {
		return fmt.Errorf("failed to delete client: %w", err)
	}

	return nil
}

// RotateSecret generates a new client secret
func (s *OIDCClientService) RotateSecret(tenantID, clientID int64) (string, error) {
	client, err := s.clientRepo.GetByID(clientID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return "", ErrClientNotFound
		}
		return "", fmt.Errorf("failed to get client: %w", err)
	}

	// Ensure client belongs to tenant
	if client.TenantID != tenantID {
		return "", ErrClientNotFound
	}

	// Cannot rotate secret for public clients
	if client.IsPublic {
		return "", errors.New("cannot rotate secret for public clients")
	}

	// Generate new secret
	clientSecret, err := generateClientSecret()
	if err != nil {
		return "", fmt.Errorf("failed to generate client secret: %w", err)
	}

	// Hash new secret
	hashedSecret, err := bcrypt.GenerateFromPassword([]byte(clientSecret), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("failed to hash client secret: %w", err)
	}

	client.ClientSecretHash = string(hashedSecret)
	if err := s.clientRepo.Update(client); err != nil {
		return "", fmt.Errorf("failed to update client: %w", err)
	}

	return clientSecret, nil
}

// toDTO converts a model to DTO
func (s *OIDCClientService) toDTO(client *models.OIDCClient) *dto.OIDCClientDTO {
	return &dto.OIDCClientDTO{
		ID:           client.ID,
		TenantID:     client.TenantID,
		ClientID:     client.ClientID,
		Name:         client.Name,
		RedirectURIs: client.RedirectURIs,
		Scopes:       client.Scopes,
		GrantTypes:   client.GrantTypes,
		IsPublic:     client.IsPublic,
		IsActive:     client.IsActive,
		CreatedAt:    client.CreatedAt.Format("2006-01-02T15:04:05Z"),
		UpdatedAt:    client.UpdatedAt.Format("2006-01-02T15:04:05Z"),
	}
}
