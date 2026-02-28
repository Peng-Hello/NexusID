package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/nexus-id/backend/internal/dto"
	"github.com/nexus-id/backend/internal/jwt"
	"github.com/nexus-id/backend/internal/service"
)

// AuthHandler handles authentication-related HTTP requests
type AuthHandler struct {
	userService   *service.UserService
	tenantService *service.TenantService
	tokenService  *jwt.TokenService
}

// NewAuthHandler creates a new AuthHandler
func NewAuthHandler(tokenService *jwt.TokenService) *AuthHandler {
	return &AuthHandler{
		userService:   service.NewUserService(),
		tenantService: service.NewTenantService(),
		tokenService:  tokenService,
	}
}

// Login handles user login
// @Summary User login
// @Tags auth
// @Accept json
// @Produce json
// @Param request body dto.LoginRequest true "Login request"
// @Success 200 {object} dto.LoginResponse
// @Router /api/v1/auth/login [post]
func (h *AuthHandler) Login(c *gin.Context) {
	var req dto.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get tenant by slug or use default
	// For now, we'll use tenant_id from request or default to 1
	// In production, you'd extract tenant from subdomain or header
	tenantID := int64(1) // Default tenant

	// Validate credentials
	user, err := h.userService.ValidatePassword(tenantID, req.Email, req.Password)
	if err != nil {
		if err == service.ErrAccountLocked {
			c.JSON(http.StatusLocked, gin.H{"error": "account is locked due to multiple failed login attempts"})
			return
		}
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid email or password"})
		return
	}

	// Get user roles
	roles := make([]string, len(user.Roles))
	for i, role := range user.Roles {
		roles[i] = role.Name
	}

	// Generate tokens
	accessToken, err := h.tokenService.GenerateAccessToken(user, "nexus-id-web", roles)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate access token"})
		return
	}

	refreshToken, err := h.tokenService.GenerateRefreshToken(user, "nexus-id-web")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate refresh token"})
		return
	}

	// Update last login
	h.userService.UpdateLastLogin(user.ID)

	c.JSON(http.StatusOK, dto.LoginResponse{
		User: dto.UserDTO{
			ID:              user.ID,
			TenantID:        user.TenantID,
			Email:           user.Email,
			FullName:        user.FullName,
			IsActive:        user.IsActive,
			IsEmailVerified: user.IsEmailVerified,
			Roles:           roles,
			CreatedAt:       user.CreatedAt.Format("2006-01-02T15:04:05Z"),
			UpdatedAt:       user.UpdatedAt.Format("2006-01-02T15:04:05Z"),
		},
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    int64(h.tokenService.GetTokenExpiry().Seconds()),
	})
}

// Register handles user registration
// @Summary User registration
// @Tags auth
// @Accept json
// @Produce json
// @Param request body dto.CreateUserRequest true "Registration request"
// @Success 200 {object} dto.UserDTO
// @Router /api/v1/auth/register [post]
func (h *AuthHandler) Register(c *gin.Context) {
	// Get tenant ID (from subdomain, header, or default)
	tenantID := int64(1) // Default tenant

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

	c.JSON(http.StatusCreated, user)
}

// RefreshToken handles token refresh
// @Summary Refresh access token
// @Tags auth
// @Accept json
// @Produce json
// @Param refresh_token body string true "Refresh token"
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/auth/refresh [post]
func (h *AuthHandler) RefreshToken(c *gin.Context) {
	var req struct {
		RefreshToken string `json:"refresh_token" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validate refresh token
	claims, err := h.tokenService.ValidateToken(req.RefreshToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid refresh token"})
		return
	}

	if claims.TokenType != jwt.RefreshTokenType {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token type"})
		return
	}

	// Get user model (not DTO)
	user, err := h.userService.GetByIDModel(claims.TenantID, claims.UserID)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not found"})
		return
	}

	// Get roles
	roles := make([]string, len(user.Roles))
	for i, role := range user.Roles {
		roles[i] = role.Name
	}

	// Generate new tokens
	accessToken, err := h.tokenService.GenerateAccessToken(user, claims.ClientID, roles)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate access token"})
		return
	}

	newRefreshToken, err := h.tokenService.GenerateRefreshToken(user, claims.ClientID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate refresh token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"access_token":  accessToken,
		"refresh_token": newRefreshToken,
		"expires_in":    int64(h.tokenService.GetTokenExpiry().Seconds()),
	})
}
