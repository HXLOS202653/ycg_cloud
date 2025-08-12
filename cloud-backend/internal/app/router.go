// Package app provides routing configuration for the application.
package app

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/HXLOS202653/ycg_cloud/cloud-backend/internal/handler"
	"github.com/HXLOS202653/ycg_cloud/cloud-backend/internal/middleware"
	"github.com/HXLOS202653/ycg_cloud/cloud-backend/internal/service"
)

// addRoutes configures all API routes for the application.
func (s *Server) addRoutes(r *gin.Engine) {
	// Initialize services
	authService := service.NewAuthService(s.Config)
	emailService := service.NewEmailService(s.Config)
	verificationService := service.NewVerificationService(s.Config, s.DB.GetRedisClient(), emailService)
	sessionService := service.NewSessionService(s.Config, s.DB.GetMySQLDB(), s.DB.GetRedisClient())
	totpService := service.NewTOTPService("YCG Cloud")
	twoFactorService := service.NewTwoFactorService(s.DB.GetMySQLDB(), s.DB.GetRedisClient(), totpService)

	// Initialize handlers
	healthHandler := handler.NewHealthHandler(s.Config, s.DB)
	authHandler := handler.NewAuthHandler(authService, emailService, verificationService, sessionService, twoFactorService, s.DB.GetMySQLDB())

	// Health check endpoints
	r.GET("/health", healthHandler.Health)
	r.GET("/health/live", healthHandler.HealthLive)
	r.GET("/health/ready", healthHandler.HealthReady)
	r.GET("/metrics", healthHandler.Metrics)
	r.GET("/", s.welcomeHandler)

	// API version 1 routes
	v1 := r.Group("/api/v1")

	// Public authentication routes (no auth required)
	auth := v1.Group("/auth")
	{
		auth.POST("/register", authHandler.Register)
		auth.POST("/login", authHandler.Login)
		auth.POST("/refresh", authHandler.RefreshToken)
		auth.POST("/forgot-password", s.placeholder("forgot password"))
		auth.POST("/reset-password", s.placeholder("reset password"))

		// Email verification routes
		auth.POST("/send-verification", authHandler.SendVerificationCode)
		auth.POST("/verify-email", authHandler.VerifyEmail)
		auth.POST("/resend-verification", authHandler.ResendVerificationCode)
	}

	// Protected routes (require authentication)
	protected := v1.Group("")
	protected.Use(middleware.AuthMiddleware(authService))
	{
		// User profile routes
		protected.GET("/auth/profile", authHandler.GetCurrentUser)
		protected.POST("/auth/logout", authHandler.Logout)
		protected.POST("/auth/logout-all", authHandler.LogoutAll)
		protected.GET("/auth/sessions", authHandler.GetActiveSessions)

		// Two-factor authentication routes
		protected.POST("/auth/2fa/setup", authHandler.Setup2FA)
		protected.POST("/auth/2fa/enable", authHandler.Enable2FA)
		protected.POST("/auth/2fa/disable", authHandler.Disable2FA)
		protected.POST("/auth/2fa/backup-codes", authHandler.RegenerateBackupCodes)
	}

	// Public 2FA verification route (no auth required for temporary verification)
	auth.POST("/2fa/verify", authHandler.Verify2FA)

	// File management routes
	files := v1.Group("/files")
	files.GET("", s.placeholder("list files"))
	files.GET("/:id", s.placeholder("get file"))
	files.POST("", s.placeholder("upload file"))
	files.PUT("/:id", s.placeholder("update file"))
	files.DELETE("/:id", s.placeholder("delete file"))
	files.POST("/:id/share", s.placeholder("share file"))
	files.GET("/:id/download", s.placeholder("download file"))
	files.GET("/:id/preview", s.placeholder("preview file"))

	// Folder management routes
	folders := v1.Group("/folders")
	folders.GET("", s.placeholder("list folders"))
	folders.GET("/:id", s.placeholder("get folder"))
	folders.POST("", s.placeholder("create folder"))
	folders.PUT("/:id", s.placeholder("update folder"))
	folders.DELETE("/:id", s.placeholder("delete folder"))

	// Search routes
	search := v1.Group("/search")
	search.GET("/files", s.placeholder("search files"))
	search.GET("/content", s.placeholder("search content"))

	// User management routes
	users := v1.Group("/users")
	users.GET("/me", s.placeholder("get current user"))
	users.PUT("/me", s.placeholder("update current user"))
	users.GET("/me/storage", s.placeholder("get storage info"))

	// Team collaboration routes
	teams := v1.Group("/teams")
	teams.GET("", s.placeholder("list teams"))
	teams.POST("", s.placeholder("create team"))
	teams.GET("/:id", s.placeholder("get team"))
	teams.PUT("/:id", s.placeholder("update team"))
	teams.DELETE("/:id", s.placeholder("delete team"))
	teams.GET("/:id/members", s.placeholder("list team members"))
	teams.POST("/:id/members", s.placeholder("add team member"))
	teams.DELETE("/:id/members/:userId", s.placeholder("remove team member"))

	// Chat routes (WebSocket will be added later)
	chat := v1.Group("/chat")
	chat.GET("/rooms", s.placeholder("list chat rooms"))
	chat.POST("/rooms", s.placeholder("create chat room"))
	chat.GET("/rooms/:id/messages", s.placeholder("get chat messages"))
	chat.POST("/rooms/:id/messages", s.placeholder("send message"))

	// Admin routes
	admin := v1.Group("/admin")
	systemHandler := handler.NewSystemHandler(s.Config, s.DB)
	admin.GET("/stats", systemHandler.SystemStats)
	admin.GET("/users", s.placeholder("list all users"))
	admin.GET("/logs", s.placeholder("system logs"))
}

// welcomeHandler handles root path requests.
func (s *Server) welcomeHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message":     "Welcome to YCG Cloud Storage API",
		"version":     s.Config.App.Version,
		"environment": s.Config.App.Environment,
		"docs":        "/api/v1/docs", // Future API documentation endpoint
	})
}

// placeholder creates a placeholder handler for unimplemented routes.
func (s *Server) placeholder(feature string) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(http.StatusNotImplemented, gin.H{
			"message":    "This endpoint is not yet implemented",
			"feature":    feature,
			"request_id": c.GetString("request_id"),
		})
	}
}
