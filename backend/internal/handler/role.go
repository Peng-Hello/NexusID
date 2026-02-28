package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/nexus-id/backend/internal/dto"
	"github.com/nexus-id/backend/internal/service"
)

// RoleHandler handles role-related HTTP requests
type RoleHandler struct {
	roleService *service.RoleService
}

// NewRoleHandler creates a new RoleHandler
func NewRoleHandler() *RoleHandler {
	return &RoleHandler{
		roleService: service.NewRoleService(),
	}
}

// Create creates a new role
// @Summary Create a new role
// @Tags roles
// @Accept json
// @Produce json
// @Param tenant_id path int true "Tenant ID"
// @Param request body dto.CreateRoleRequest true "Create role request"
// @Success 200 {object} dto.RoleDTO
// @Router /api/v1/tenants/{tenant_id}/roles [post]
func (h *RoleHandler) Create(c *gin.Context) {
	tenantID, err := strconv.ParseInt(c.Param("tenant_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid tenant ID"})
		return
	}

	var req dto.CreateRoleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	role, err := h.roleService.Create(tenantID, &req)
	if err != nil {
		if err == service.ErrTenantNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "tenant not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, role)
}

// GetByID retrieves a role by ID
// @Summary Get role by ID
// @Tags roles
// @Accept json
// @Produce json
// @Param tenant_id path int true "Tenant ID"
// @Param id path int true "Role ID"
// @Success 200 {object} dto.RoleDTO
// @Router /api/v1/tenants/{tenant_id}/roles/{id} [get]
func (h *RoleHandler) GetByID(c *gin.Context) {
	tenantID, err := strconv.ParseInt(c.Param("tenant_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid tenant ID"})
		return
	}

	roleID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid role ID"})
		return
	}

	role, err := h.roleService.GetByID(tenantID, roleID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "role not found"})
		return
	}

	c.JSON(http.StatusOK, role)
}

// List retrieves roles for a tenant with pagination
// @Summary List roles
// @Tags roles
// @Accept json
// @Produce json
// @Param tenant_id path int true "Tenant ID"
// @Param page query int false "Page number" default(1)
// @Param page_size query int false "Page size" default(20)
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/tenants/{tenant_id}/roles [get]
func (h *RoleHandler) List(c *gin.Context) {
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

	roles, total, err := h.roleService.List(tenantID, page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list roles"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":      roles,
		"total":     total,
		"page":      page,
		"page_size": pageSize,
	})
}

// Update updates a role
// @Summary Update role
// @Tags roles
// @Accept json
// @Produce json
// @Param tenant_id path int true "Tenant ID"
// @Param id path int true "Role ID"
// @Param request body dto.UpdateRoleRequest true "Update role request"
// @Success 200 {object} dto.RoleDTO
// @Router /api/v1/tenants/{tenant_id}/roles/{id} [put]
func (h *RoleHandler) Update(c *gin.Context) {
	tenantID, err := strconv.ParseInt(c.Param("tenant_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid tenant ID"})
		return
	}

	roleID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid role ID"})
		return
	}

	var req dto.UpdateRoleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	role, err := h.roleService.Update(tenantID, roleID, &req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, role)
}

// Delete deletes a role
// @Summary Delete role
// @Tags roles
// @Accept json
// @Produce json
// @Param tenant_id path int true "Tenant ID"
// @Param id path int true "Role ID"
// @Success 200 {object} map[string]string
// @Router /api/v1/tenants/{tenant_id}/roles/{id} [delete]
func (h *RoleHandler) Delete(c *gin.Context) {
	tenantID, err := strconv.ParseInt(c.Param("tenant_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid tenant ID"})
		return
	}

	roleID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid role ID"})
		return
	}

	if err := h.roleService.Delete(tenantID, roleID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "role deleted successfully"})
}

// AssignRole assigns a role to a user
// @Summary Assign role to user
// @Tags roles
// @Accept json
// @Produce json
// @Param tenant_id path int true "Tenant ID"
// @Param request body dto.AssignRoleRequest true "Assign role request"
// @Success 200 {object} map[string]string
// @Router /api/v1/tenants/{tenant_id}/roles/assign [post]
func (h *RoleHandler) AssignRole(c *gin.Context) {
	tenantID, err := strconv.ParseInt(c.Param("tenant_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid tenant ID"})
		return
	}

	var req dto.AssignRoleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.roleService.AssignRole(tenantID, req.UserID, req.RoleID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "role assigned successfully"})
}

// RevokeRole revokes a role from a user
// @Summary Revoke role from user
// @Tags roles
// @Accept json
// @Produce json
// @Param tenant_id path int true "Tenant ID"
// @Param request body dto.AssignRoleRequest true "Revoke role request"
// @Success 200 {object} map[string]string
// @Router /api/v1/tenants/{tenant_id}/roles/revoke [post]
func (h *RoleHandler) RevokeRole(c *gin.Context) {
	tenantID, err := strconv.ParseInt(c.Param("tenant_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid tenant ID"})
		return
	}

	var req dto.AssignRoleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.roleService.RevokeRole(tenantID, req.UserID, req.RoleID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "role revoked successfully"})
}

// GetUserRoles retrieves all roles for a user
// @Summary Get user roles
// @Tags roles
// @Accept json
// @Produce json
// @Param tenant_id path int true "Tenant ID"
// @Param user_id path int true "User ID"
// @Success 200 {object} []dto.RoleDTO
// @Router /api/v1/tenants/{tenant_id}/users/{user_id}/roles [get]
func (h *RoleHandler) GetUserRoles(c *gin.Context) {
	userID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user ID"})
		return
	}

	roles, err := h.roleService.GetUserRoles(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get user roles"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": roles})
}
