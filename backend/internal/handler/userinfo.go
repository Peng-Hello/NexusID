package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	nexusjwt "github.com/nexus-id/backend/internal/jwt"
	"github.com/nexus-id/backend/internal/service"
)

// UserInfoHandler handles the OIDC userinfo endpoint
type UserInfoHandler struct {
	userService *service.UserService
}

// NewUserInfoHandler creates a new UserInfoHandler
func NewUserInfoHandler() *UserInfoHandler {
	return &UserInfoHandler{
		userService: service.NewUserService(),
	}
}

// GetUserInfo returns the authenticated user's claims (OIDC userinfo endpoint)
// @Summary OIDC UserInfo endpoint
// @Tags oauth
// @Security BearerAuth
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Router /oauth/userinfo [get]
func (h *UserInfoHandler) GetUserInfo(c *gin.Context) {
	// Claims were set by AuthRequired middleware
	claims, exists := c.Get("claims")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	tokenClaims, ok := claims.(*nexusjwt.TokenClaims)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid claims"})
		return
	}

	userID := tokenClaims.UserID
	tenantID := tokenClaims.TenantID

	// Fetch fresh user info from DB
	user, err := h.userService.GetByID(tenantID, userID)
	if err != nil {
		// Fall back to claims in JWT if DB lookup fails
		c.JSON(http.StatusOK, gin.H{
			"sub":            tokenClaims.Subject,
			"email":          tokenClaims.Email,
			"name":           tokenClaims.FullName,
			"email_verified": false,
			"tenant_id":      tenantID,
			"roles":          tokenClaims.Roles,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"sub":            tokenClaims.Subject,
		"email":          user.Email,
		"name":           user.FullName,
		"email_verified": user.IsEmailVerified,
		"tenant_id":      user.TenantID,
		"roles":          user.Roles,
		"updated_at":     user.UpdatedAt,
	})
}
