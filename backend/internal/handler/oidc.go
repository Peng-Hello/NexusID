package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/nexus-id/backend/internal/jwt"
)

// OIDCHandler handles OIDC protocol endpoints
type OIDCHandler struct {
	keyManager *jwt.KeyManager
	issuerURL  string
}

// NewOIDCHandler creates a new OIDC handler
func NewOIDCHandler(keyManager *jwt.KeyManager, issuerURL string) *OIDCHandler {
	return &OIDCHandler{
		keyManager: keyManager,
		issuerURL:  issuerURL,
	}
}

// Discovery returns the OIDC discovery document
func (h *OIDCHandler) Discovery(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"issuer":                                h.issuerURL,
		"authorization_endpoint":                h.issuerURL + "/oauth/authorize",
		"token_endpoint":                        h.issuerURL + "/oauth/token",
		"jwks_uri":                              h.issuerURL + "/.well-known/jwks.json",
		"response_types_supported":              []string{"code"},
		"subject_types_supported":               []string{"public"},
		"id_token_signing_alg_values_supported": []string{"RS256"},
		"scopes_supported":                      []string{"openid", "profile", "email"},
		"grant_types_supported":                 []string{"authorization_code", "refresh_token"},
		"token_endpoint_auth_methods_supported": []string{"client_secret_basic", "client_secret_post"},
		"claims_supported":                      []string{"sub", "name", "email", "email_verified", "preferred_username"},
	})
}

// JWKS returns the JSON Web Key Set
func (h *OIDCHandler) JWKS(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"keys": []interface{}{h.keyManager.GetPublicJWK()},
	})
}
