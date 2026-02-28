package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/nexus-id/backend/internal/dto"
	"github.com/nexus-id/backend/internal/service"
)

// OIDCClientHandler handles OIDC client-related HTTP requests
type OIDCClientHandler struct {
	clientService *service.OIDCClientService
}

// NewOIDCClientHandler creates a new OIDCClientHandler
func NewOIDCClientHandler() *OIDCClientHandler {
	return &OIDCClientHandler{
		clientService: service.NewOIDCClientService(),
	}
}

// Create creates a new OIDC client
// @Summary Create a new OIDC client
// @Tags oidc-clients
// @Accept json
// @Produce json
// @Param tenant_id path int true "Tenant ID"
// @Param request body dto.CreateOIDCClientRequest true "Create client request"
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/tenants/{tenant_id}/clients [post]
func (h *OIDCClientHandler) Create(c *gin.Context) {
	tenantID, err := strconv.ParseInt(c.Param("tenant_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid tenant ID"})
		return
	}

	var req dto.CreateOIDCClientRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	client, secret, err := h.clientService.Create(tenantID, &req)
	if err != nil {
		if err == service.ErrTenantNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "tenant not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create client"})
		return
	}

	// Return client with secret (only shown once)
	c.JSON(http.StatusOK, gin.H{
		"client":  client,
		"secret":  secret,
		"message": "Save the client secret securely. It will not be shown again.",
	})
}

// GetByID retrieves an OIDC client by ID
// @Summary Get OIDC client by ID
// @Tags oidc-clients
// @Accept json
// @Produce json
// @Param tenant_id path int true "Tenant ID"
// @Param id path int true "Client ID"
// @Success 200 {object} dto.OIDCClientDTO
// @Router /api/v1/tenants/{tenant_id}/clients/{id} [get]
func (h *OIDCClientHandler) GetByID(c *gin.Context) {
	tenantID, err := strconv.ParseInt(c.Param("tenant_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid tenant ID"})
		return
	}

	clientID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid client ID"})
		return
	}

	client, err := h.clientService.GetByID(tenantID, clientID)
	if err != nil {
		if err == service.ErrClientNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "client not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get client"})
		return
	}

	c.JSON(http.StatusOK, client)
}

// List retrieves OIDC clients for a tenant with pagination
// @Summary List OIDC clients
// @Tags oidc-clients
// @Accept json
// @Produce json
// @Param tenant_id path int true "Tenant ID"
// @Param page query int false "Page number" default(1)
// @Param page_size query int false "Page size" default(20)
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/tenants/{tenant_id}/clients [get]
func (h *OIDCClientHandler) List(c *gin.Context) {
	tenantID, err := strconv.ParseInt(c.Param("tenant_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid tenant ID"})
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	clients, total, err := h.clientService.List(tenantID, page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list clients"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":      clients,
		"total":     total,
		"page":      page,
		"page_size": pageSize,
	})
}

// Update updates an OIDC client
// @Summary Update OIDC client
// @Tags oidc-clients
// @Accept json
// @Produce json
// @Param tenant_id path int true "Tenant ID"
// @Param id path int true "Client ID"
// @Param request body dto.UpdateOIDCClientRequest true "Update client request"
// @Success 200 {object} dto.OIDCClientDTO
// @Router /api/v1/tenants/{tenant_id}/clients/{id} [put]
func (h *OIDCClientHandler) Update(c *gin.Context) {
	tenantID, err := strconv.ParseInt(c.Param("tenant_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid tenant ID"})
		return
	}

	clientID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid client ID"})
		return
	}

	var req dto.UpdateOIDCClientRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	client, err := h.clientService.Update(tenantID, clientID, &req)
	if err != nil {
		if err == service.ErrClientNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "client not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update client"})
		return
	}

	c.JSON(http.StatusOK, client)
}

// Delete soft deletes an OIDC client
// @Summary Delete OIDC client
// @Tags oidc-clients
// @Accept json
// @Produce json
// @Param tenant_id path int true "Tenant ID"
// @Param id path int true "Client ID"
// @Success 200 {object} map[string]string
// @Router /api/v1/tenants/{tenant_id}/clients/{id} [delete]
func (h *OIDCClientHandler) Delete(c *gin.Context) {
	tenantID, err := strconv.ParseInt(c.Param("tenant_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid tenant ID"})
		return
	}

	clientID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid client ID"})
		return
	}

	if err := h.clientService.Delete(tenantID, clientID); err != nil {
		if err == service.ErrClientNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "client not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete client"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "client deleted successfully"})
}

// RotateSecret generates a new client secret
// @Summary Rotate client secret
// @Tags oidc-clients
// @Accept json
// @Produce json
// @Param tenant_id path int true "Tenant ID"
// @Param id path int true "Client ID"
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/tenants/{tenant_id}/clients/{id}/rotate-secret [post]
func (h *OIDCClientHandler) RotateSecret(c *gin.Context) {
	tenantID, err := strconv.ParseInt(c.Param("tenant_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid tenant ID"})
		return
	}

	clientID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid client ID"})
		return
	}

	secret, err := h.clientService.RotateSecret(tenantID, clientID)
	if err != nil {
		if err == service.ErrClientNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "client not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"secret":  secret,
		"message": "Save the new client secret securely. It will not be shown again.",
	})
}
