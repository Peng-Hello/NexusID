package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/nexus-id/backend/internal/dto"
	"github.com/nexus-id/backend/internal/service"
)

// TenantHandler handles tenant-related HTTP requests
type TenantHandler struct {
	tenantService *service.TenantService
}

// NewTenantHandler creates a new TenantHandler
func NewTenantHandler() *TenantHandler {
	return &TenantHandler{
		tenantService: service.NewTenantService(),
	}
}

// Create creates a new tenant
// @Summary Create a new tenant
// @Tags tenants
// @Accept json
// @Produce json
// @Param request body dto.CreateTenantRequest true "Create tenant request"
// @Success 200 {object} dto.TenantDTO
// @Router /api/v1/tenants [post]
func (h *TenantHandler) Create(c *gin.Context) {
	var req dto.CreateTenantRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	tenant, err := h.tenantService.Create(&req)
	if err != nil {
		if err == service.ErrTenantSlugExists {
			c.JSON(http.StatusConflict, gin.H{"error": "tenant slug already exists"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create tenant"})
		return
	}

	c.JSON(http.StatusOK, tenant)
}

// GetByID retrieves a tenant by ID
// @Summary Get tenant by ID
// @Tags tenants
// @Accept json
// @Produce json
// @Param id path int true "Tenant ID"
// @Success 200 {object} dto.TenantDTO
// @Router /api/v1/tenants/{id} [get]
func (h *TenantHandler) GetByID(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("tenant_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid tenant ID"})
		return
	}

	tenant, err := h.tenantService.GetByID(id)
	if err != nil {
		if err == service.ErrTenantNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "tenant not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get tenant"})
		return
	}

	c.JSON(http.StatusOK, tenant)
}

// GetBySlug retrieves a tenant by slug
// @Summary Get tenant by slug
// @Tags tenants
// @Accept json
// @Produce json
// @Param slug path string true "Tenant slug"
// @Success 200 {object} dto.TenantDTO
// @Router /api/v1/tenants/slug/{slug} [get]
func (h *TenantHandler) GetBySlug(c *gin.Context) {
	slug := c.Param("slug")

	tenant, err := h.tenantService.GetBySlug(slug)
	if err != nil {
		if err == service.ErrTenantNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "tenant not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get tenant"})
		return
	}

	c.JSON(http.StatusOK, tenant)
}

// List retrieves tenants with pagination
// @Summary List tenants
// @Tags tenants
// @Accept json
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param page_size query int false "Page size" default(20)
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/tenants [get]
func (h *TenantHandler) List(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	tenants, total, err := h.tenantService.List(page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list tenants"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":      tenants,
		"total":     total,
		"page":      page,
		"page_size": pageSize,
	})
}

// Update updates a tenant
// @Summary Update tenant
// @Tags tenants
// @Accept json
// @Produce json
// @Param id path int true "Tenant ID"
// @Param request body dto.UpdateTenantRequest true "Update tenant request"
// @Success 200 {object} dto.TenantDTO
// @Router /api/v1/tenants/{id} [put]
func (h *TenantHandler) Update(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("tenant_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid tenant ID"})
		return
	}

	var req dto.UpdateTenantRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	tenant, err := h.tenantService.Update(id, &req)
	if err != nil {
		if err == service.ErrTenantNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "tenant not found"})
			return
		}
		if err == service.ErrTenantSlugExists {
			c.JSON(http.StatusConflict, gin.H{"error": "tenant slug already exists"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update tenant"})
		return
	}

	c.JSON(http.StatusOK, tenant)
}

// Delete soft deletes a tenant
// @Summary Delete tenant
// @Tags tenants
// @Accept json
// @Produce json
// @Param id path int true "Tenant ID"
// @Success 200 {object} map[string]string
// @Router /api/v1/tenants/{id} [delete]
func (h *TenantHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("tenant_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid tenant ID"})
		return
	}

	if err := h.tenantService.Delete(id); err != nil {
		if err == service.ErrTenantNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "tenant not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete tenant"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "tenant deleted successfully"})
}

// UserHandler handles user-related HTTP requests
type UserHandler struct {
	userService *service.UserService
}

// NewUserHandler creates a new UserHandler
func NewUserHandler() *UserHandler {
	return &UserHandler{
		userService: service.NewUserService(),
	}
}

// Create creates a new user
// @Summary Create a new user
// @Tags users
// @Accept json
// @Produce json
// @Param tenant_id path int true "Tenant ID"
// @Param request body dto.CreateUserRequest true "Create user request"
// @Success 200 {object} dto.UserDTO
// @Router /api/v1/tenants/{tenant_id}/users [post]
func (h *UserHandler) Create(c *gin.Context) {
	tenantID, err := strconv.ParseInt(c.Param("tenant_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid tenant ID"})
		return
	}

	var req dto.CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := h.userService.Create(tenantID, &req)
	if err != nil {
		if err == service.ErrTenantNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "tenant not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create user"})
		return
	}

	c.JSON(http.StatusOK, user)
}

// GetByID retrieves a user by ID
// @Summary Get user by ID
// @Tags users
// @Accept json
// @Produce json
// @Param tenant_id path int true "Tenant ID"
// @Param id path int true "User ID"
// @Success 200 {object} dto.UserDTO
// @Router /api/v1/tenants/{tenant_id}/users/{id} [get]
func (h *UserHandler) GetByID(c *gin.Context) {
	tenantID, err := strconv.ParseInt(c.Param("tenant_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid tenant ID"})
		return
	}

	userID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user ID"})
		return
	}

	user, err := h.userService.GetByID(tenantID, userID)
	if err != nil {
		if err == service.ErrUserNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get user"})
		return
	}

	c.JSON(http.StatusOK, user)
}

// List retrieves users for a tenant with pagination
// @Summary List users
// @Tags users
// @Accept json
// @Produce json
// @Param tenant_id path int true "Tenant ID"
// @Param page query int false "Page number" default(1)
// @Param page_size query int false "Page size" default(20)
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/tenants/{tenant_id}/users [get]
func (h *UserHandler) List(c *gin.Context) {
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

	users, total, err := h.userService.List(tenantID, page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list users"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":      users,
		"total":     total,
		"page":      page,
		"page_size": pageSize,
	})
}

// Update updates a user
// @Summary Update user
// @Tags users
// @Accept json
// @Produce json
// @Param tenant_id path int true "Tenant ID"
// @Param id path int true "User ID"
// @Param request body dto.UpdateUserRequest true "Update user request"
// @Success 200 {object} dto.UserDTO
// @Router /api/v1/tenants/{tenant_id}/users/{id} [put]
func (h *UserHandler) Update(c *gin.Context) {
	tenantID, err := strconv.ParseInt(c.Param("tenant_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid tenant ID"})
		return
	}

	userID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user ID"})
		return
	}

	var req dto.UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := h.userService.Update(tenantID, userID, &req)
	if err != nil {
		if err == service.ErrUserNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update user"})
		return
	}

	c.JSON(http.StatusOK, user)
}

// Delete soft deletes a user
// @Summary Delete user
// @Tags users
// @Accept json
// @Produce json
// @Param tenant_id path int true "Tenant ID"
// @Param id path int true "User ID"
// @Success 200 {object} map[string]string
// @Router /api/v1/tenants/{tenant_id}/users/{id} [delete]
func (h *UserHandler) Delete(c *gin.Context) {
	tenantID, err := strconv.ParseInt(c.Param("tenant_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid tenant ID"})
		return
	}

	userID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user ID"})
		return
	}

	if err := h.userService.Delete(tenantID, userID); err != nil {
		if err == service.ErrUserNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "user deleted successfully"})
}

// ChangePassword handles password change requests
// @Summary Change user password
// @Tags users
// @Accept json
// @Produce json
// @Param request body dto.ChangePasswordRequest true "Change password request"
// @Security BearerAuth
// @Success 200 {object} map[string]string
// @Router /api/v1/user/change-password [post]
func (h *UserHandler) ChangePassword(c *gin.Context) {
	var req dto.ChangePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get user info from JWT context
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	tenantID, exists := c.Get("tenant_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	// Call service
	err := h.userService.ChangePassword(
		tenantID.(int64),
		userID.(int64),
		req.NewPassword,
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to change password"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "password changed successfully"})
}
