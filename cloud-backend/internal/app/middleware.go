// Package app provides middleware configuration for the application.
package app

import (
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// addMiddlewares adds all necessary middlewares to the Gin router.
func (s *Server) addMiddlewares(r *gin.Engine) {
	// Recovery middleware
	r.Use(gin.Recovery())

	// Custom logger middleware
	r.Use(loggerMiddleware())

	// CORS middleware
	r.Use(corsMiddleware())

	// Request ID middleware
	r.Use(requestIDMiddleware())

	// Security headers middleware
	r.Use(securityHeadersMiddleware())

	// Rate limiting middleware (if needed)
	// r.Use(rateLimitMiddleware())
}

// loggerMiddleware creates a structured logger middleware using logrus.
func loggerMiddleware() gin.HandlerFunc {
	logger := logrus.New()
	logger.SetFormatter(&logrus.JSONFormatter{})

	return gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
		logger.WithFields(logrus.Fields{
			"status_code": param.StatusCode,
			"latency":     param.Latency,
			"client_ip":   param.ClientIP,
			"method":      param.Method,
			"path":        param.Path,
			"error":       param.ErrorMessage,
			"user_agent":  param.Request.UserAgent(),
			"timestamp":   param.TimeStamp.Format(time.RFC3339),
		}).Info("HTTP Request")
		return ""
	})
}

// corsMiddleware configures CORS (Cross-Origin Resource Sharing) settings.
func corsMiddleware() gin.HandlerFunc {
	config := cors.Config{
		AllowOrigins:     []string{"http://localhost:3000", "http://localhost:5173"}, // Add your frontend URLs
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization", "X-Requested-With", "X-Request-ID"},
		ExposeHeaders:    []string{"Content-Length", "X-Request-ID"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}
	return cors.New(config)
}

// requestIDMiddleware adds a unique request ID to each request.
func requestIDMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Check if request ID already exists in header
		requestID := c.GetHeader("X-Request-ID")
		if requestID == "" {
			// Generate a simple request ID (in production, consider using UUID)
			requestID = generateRequestID()
		}

		// Set request ID in context and response header
		c.Set("request_id", requestID)
		c.Header("X-Request-ID", requestID)

		c.Next()
	}
}

// securityHeadersMiddleware adds security headers to responses.
func securityHeadersMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Security headers
		c.Header("X-Content-Type-Options", "nosniff")
		c.Header("X-Frame-Options", "DENY")
		c.Header("X-XSS-Protection", "1; mode=block")
		c.Header("Referrer-Policy", "strict-origin-when-cross-origin")

		// Remove server information
		c.Header("Server", "")

		c.Next()
	}
}

// generateRequestID generates a simple request ID based on timestamp.
// In production, consider using a more robust UUID library.
func generateRequestID() string {
	return time.Now().Format("20060102150405.000000")
}
