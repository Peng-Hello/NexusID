package main

import (
	"flag"
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/nexus-id/backend/internal/cache"
	"github.com/nexus-id/backend/internal/config"
	"github.com/nexus-id/backend/internal/database"
	"github.com/nexus-id/backend/internal/handler"
	"github.com/nexus-id/backend/internal/jwt"
	"github.com/nexus-id/backend/internal/logger"
	"github.com/nexus-id/backend/internal/middleware"
	"go.uber.org/zap"
)

var (
	configPath string
)

func init() {
	flag.StringVar(&configPath, "config", "", "Path to config file")
}

func main() {
	flag.Parse()

	// Load configuration
	cfg, err := config.LoadConfig(configPath)
	if err != nil {
		panic("Failed to load config: " + err.Error())
	}

	// Initialize logger
	if err := logger.Init(cfg.Server.Mode); err != nil {
		panic("Failed to initialize logger: " + err.Error())
	}
	defer logger.Sync()
	log := logger.Get()

	log.Info("Starting NexusID",
		zap.String("mode", cfg.Server.Mode),
		zap.Int("port", cfg.Server.Port),
	)

	// Initialize database
	if err := database.Init(cfg, log); err != nil {
		panic("Failed to initialize database: " + err.Error())
	}

	// Initialize cache
	cache, err := cache.New(cfg, log)
	if err != nil {
		panic("Failed to initialize cache: " + err.Error())
	}
	defer cache.Close()

	// Initialize JWT key manager
	keyManager, err := jwt.NewKeyManager(&cfg.JWT, log)
	if err != nil {
		panic("Failed to initialize key manager: " + err.Error())
	}

	// Initialize token service
	tokenService := jwt.NewTokenService(keyManager, cfg.JWT.Issuer, cfg.JWT.AccessTokenExpiry, cfg.JWT.RefreshTokenExpiry)

	// Initialize rate limiter
	rateLimiter := middleware.NewRateLimiter(cache, log)

	// Set Gin mode
	gin.SetMode(cfg.Server.Mode)

	// Initialize Gin router
	r := gin.New()

	// Recovery and custom logging middleware
	r.Use(gin.Recovery())
	r.Use(loggingMiddleware(log))

	// Initialize handlers
	healthHandler := handler.NewHealthHandler()
	oidcHandler := handler.NewOIDCHandler(keyManager, cfg.JWT.Issuer)
	tenantHandler := handler.NewTenantHandler()
	userHandler := handler.NewUserHandler()
	clientHandler := handler.NewOIDCClientHandler()
	roleHandler := handler.NewRoleHandler()
	authHandler := handler.NewAuthHandler(tokenService)
	oauthFlowHandler := handler.NewOIDCFlowHandler(tokenService)

	// Auth middleware
	authMiddleware := middleware.AuthRequired(tokenService)

	// API v1 routes
	v1 := r.Group("/api/v1")
	{
		// Auth routes
		auth := v1.Group("/auth")
		{
			auth.POST("/login", middleware.LoginRateLimitMiddleware(rateLimiter), authHandler.Login)
			auth.POST("/register", authHandler.Register)
			auth.POST("/refresh", authHandler.RefreshToken)
		}

		// Tenant routes
		tenants := v1.Group("/tenants")
		{
			tenants.POST("", tenantHandler.Create)
			tenants.GET("", tenantHandler.List)
			tenants.GET("/slug/:slug", tenantHandler.GetBySlug)
			tenants.GET("/:tenant_id", tenantHandler.GetByID)
			tenants.PUT("/:tenant_id", tenantHandler.Update)
			tenants.DELETE("/:tenant_id", tenantHandler.Delete)

			// User routes within tenant
			tenants.GET("/:tenant_id/users", userHandler.List)
			tenants.POST("/:tenant_id/users", userHandler.Create)
			tenants.GET("/:tenant_id/users/:id", userHandler.GetByID)
			tenants.PUT("/:tenant_id/users/:id", userHandler.Update)
			tenants.DELETE("/:tenant_id/users/:id", userHandler.Delete)
			tenants.GET("/:tenant_id/users/:id/roles", roleHandler.GetUserRoles)

			// Role routes within tenant
			tenants.GET("/:tenant_id/roles", roleHandler.List)
			tenants.POST("/:tenant_id/roles", roleHandler.Create)
			tenants.GET("/:tenant_id/roles/:id", roleHandler.GetByID)
			tenants.PUT("/:tenant_id/roles/:id", roleHandler.Update)
			tenants.DELETE("/:tenant_id/roles/:id", roleHandler.Delete)
			tenants.POST("/:tenant_id/roles/assign", roleHandler.AssignRole)
			tenants.POST("/:tenant_id/roles/revoke", roleHandler.RevokeRole)

			// OIDC client routes within tenant
			tenants.GET("/:tenant_id/clients", clientHandler.List)
			tenants.POST("/:tenant_id/clients", clientHandler.Create)
			tenants.GET("/:tenant_id/clients/:id", clientHandler.GetByID)
			tenants.PUT("/:tenant_id/clients/:id", clientHandler.Update)
			tenants.DELETE("/:tenant_id/clients/:id", clientHandler.Delete)
			tenants.POST("/:tenant_id/clients/:id/rotate-secret", clientHandler.RotateSecret)
		}
	}

	// Health check endpoints
	r.GET("/health", healthHandler.Health)
	r.GET("/ready", healthHandler.Ready)

	// Protected OIDC endpoints (requires authentication)
	protectedOauth := v1.Group("/oauth")
	protectedOauth.Use(authMiddleware)
	{
		protectedOauth.POST("/consent", oauthFlowHandler.Consent)
		protectedOauth.POST("/revoke", oauthFlowHandler.Revoke)
	}

	// OIDC endpoints
	r.GET("/.well-known/openid-configuration", oidcHandler.Discovery)
	r.GET("/.well-known/jwks.json", oidcHandler.JWKS)

	// OAuth2/OIDC flow endpoints
	r.GET("/oauth/authorize", oauthFlowHandler.Authorize)
	r.POST("/oauth/token", oauthFlowHandler.Token)
	r.POST("/oauth/revoke", oauthFlowHandler.Revoke)
	r.GET("/oauth/logout", oauthFlowHandler.Logout)

	// Start server
	addr := fmt.Sprintf(":%d", cfg.Server.Port)
	log.Info("Server started", zap.String("addr", addr))

	if err := r.Run(addr); err != nil {
		log.Fatal("Failed to start server", zap.Error(err))
	}
}

func loggingMiddleware(log *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		path := c.Request.URL.Path
		start := time.Now()

		c.Next()

		latency := time.Since(start)
		statusCode := c.Writer.Status()

		log.Info("HTTP Request",
			zap.String("method", c.Request.Method),
			zap.String("path", path),
			zap.Int("status", statusCode),
			zap.Duration("latency", latency),
			zap.String("ip", c.ClientIP()),
		)
	}
}
