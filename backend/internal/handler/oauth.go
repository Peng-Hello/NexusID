package handler

import (
	"net/http"
	"net/url"

	"github.com/gin-gonic/gin"
	nexusjwt "github.com/nexus-id/backend/internal/jwt"
	"github.com/nexus-id/backend/internal/service"
)

// OIDCFlowHandler handles OIDC flow endpoints
type OIDCFlowHandler struct {
	oidcService *service.OIDCService
}

// NewOIDCFlowHandler creates a new OIDCFlowHandler
func NewOIDCFlowHandler(tokenService *nexusjwt.TokenService) *OIDCFlowHandler {
	return &OIDCFlowHandler{
		oidcService: service.NewOIDCService(tokenService),
	}
}

// Authorize handles the authorization endpoint
// @Summary OAuth2 Authorization endpoint
// @Tags oauth
// @Accept json
// @Produce json
// @Param response_type query string true "Response type (must be 'code')"
// @Param client_id query string true "Client ID"
// @Param redirect_uri query string true "Redirect URI"
// @Param scope query string false "Scope"
// @Param state query string false "State"
// @Success 307 {object} nil
// @Router /oauth/authorize [get]
func (h *OIDCFlowHandler) Authorize(c *gin.Context) {
	var req service.AuthorizeRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.HTML(http.StatusBadRequest, "error.html", gin.H{
			"error":             "invalid_request",
			"error_description": err.Error(),
		})
		return
	}

	// For NexusID, the OAuth authorization process happens in the frontend
	// We redirect the user to the frontend consent screen, appending the query params.
	// The frontend will check if the user is logged in, and if so, call /api/v1/oauth/consent.

	// Create frontend consent URL
	consentURL, err := url.Parse("http://localhost:5173/oauth/consent")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate consent URL"})
		return
	}

	// Transfer query parameters
	values := consentURL.Query()
	values.Add("client_id", req.ClientID)
	values.Add("redirect_uri", req.RedirectURI)
	values.Add("response_type", req.ResponseType)

	if req.Scope != "" {
		values.Add("scope", req.Scope)
	}
	if req.State != "" {
		values.Add("state", req.State)
	}

	consentURL.RawQuery = values.Encode()
	c.Redirect(http.StatusTemporaryRedirect, consentURL.String())
}

// Consent handles user approval from the frontend consent screen
// @Summary Approve OAuth2 Authorization
// @Tags oauth
// @Accept json
// @Produce json
// @Param request body service.AuthorizeRequest true "Consent approval"
// @Success 200 {object} map[string]string
// @Router /api/v1/oauth/consent [post]
func (h *OIDCFlowHandler) Consent(c *gin.Context) {
	// User is authenticated via AuthMiddleware
	claims, exists := c.Get("claims")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	userClaims := claims.(*nexusjwt.TokenClaims)

	var req service.AuthorizeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":             "invalid_request",
			"error_description": err.Error(),
		})
		return
	}

	// Set the user identity from the claims
	req.UserID = userClaims.UserID
	req.TenantID = userClaims.TenantID

	// Generate authorization code
	code, err := h.oidcService.GenerateAuthorizationCode(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":             "invalid_request",
			"error_description": err.Error(),
		})
		return
	}

	// Build the redirect URL for the client application
	redirectURL, err := url.Parse(req.RedirectURI)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":             "invalid_request",
			"error_description": "invalid redirect_uri format",
		})
		return
	}

	params := redirectURL.Query()
	params.Set("code", code)
	if req.State != "" {
		params.Set("state", req.State)
	}
	redirectURL.RawQuery = params.Encode()

	// Return the redirect URL to the frontend so it can perform the redirect
	c.JSON(http.StatusOK, gin.H{
		"redirect_url": redirectURL.String(),
	})
}

// Token handles the token endpoint
// @Summary OAuth2 Token endpoint
// @Tags oauth
// @Accept json
// @Produce json
// @Param grant_type formData string true "Grant type"
// @Param code formData string false "Authorization code"
// @Param redirect_uri formData string false "Redirect URI"
// @Param client_id formData string true "Client ID"
// @Param client_secret formData string false "Client secret"
// @Param refresh_token formData string false "Refresh token"
// @Success 200 {object} service.TokenResponse
// @Router /oauth/token [post]
func (h *OIDCFlowHandler) Token(c *gin.Context) {
	var req service.TokenRequest

	// Try to bind from form data (for POST with application/x-www-form-urlencoded)
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":             "invalid_request",
			"error_description": err.Error(),
		})
		return
	}

	var response *service.TokenResponse
	var err error

	switch req.GrantType {
	case "authorization_code":
		response, err = h.oidcService.ExchangeCodeForToken(&req)
	case "refresh_token":
		response, err = h.oidcService.RefreshAccessToken(&req)
	default:
		c.JSON(http.StatusBadRequest, gin.H{
			"error":             "unsupported_grant_type",
			"error_description": "Grant type not supported",
		})
		return
	}

	if err != nil {
		statusCode := http.StatusBadRequest
		if err == service.ErrInvalidClient {
			statusCode = http.StatusUnauthorized
		}
		c.JSON(statusCode, gin.H{
			"error":             "invalid_grant",
			"error_description": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, response)
}

// Revoke handles token revocation (for logout)
// @Summary Revoke token
// @Tags oauth
// @Accept json
// @Produce json
// @Param token body string true "Token to revoke"
// @Param token_type_hint body string false "Token type hint"
// @Success 200 {object} map[string]string
// @Router /oauth/revoke [post]
func (h *OIDCFlowHandler) Revoke(c *gin.Context) {
	// Parse the request using simple form binding
	var req struct {
		Token         string `form:"token" json:"token" binding:"required"`
		TokenTypeHint string `form:"token_type_hint" json:"token_type_hint"`
	}

	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid_request"})
		return
	}

	err := h.oidcService.RevokeToken(req.Token, req.TokenTypeHint)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"status": "ignored or failed, but returning 200 per spec"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "revoked"})
}

// Logout handles the end session endpoint
// @Summary End session
// @Tags oauth
// @Accept json
// @Produce json
// @Param id_token_hint query string false "ID token hint"
// @Param post_logout_redirect_uri query string false "Post logout redirect URI"
// @Param state query string false "State"
// @Success 200 {object} map[string]string
// @Router /oauth/logout [get]
func (h *OIDCFlowHandler) Logout(c *gin.Context) {
	idTokenHint := c.Query("id_token_hint")
	postLogoutRedirectURI := c.Query("post_logout_redirect_uri")
	state := c.Query("state")

	c.SetCookie("auth_token", "", -1, "/", "", false, true)

	_ = idTokenHint
	_ = state

	if postLogoutRedirectURI != "" {
		redirectURL, err := url.Parse(postLogoutRedirectURI)
		if err == nil {
			if state != "" {
				params := redirectURL.Query()
				params.Set("state", state)
				redirectURL.RawQuery = params.Encode()
			}
			c.Redirect(http.StatusFound, redirectURL.String())
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{"message": "logged out successfully"})
}
