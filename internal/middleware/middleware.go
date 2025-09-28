package middleware

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// Logger records HTTP request metadata
func Logger(log *logrus.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Capture start time
		startTime := time.Now()

		// Process request
		c.Next()

		// Capture end time
		endTime := time.Now()
		// Compute latency
		latencyTime := endTime.Sub(startTime)

		// Request method
		reqMethod := c.Request.Method
		// Request URI
		reqURI := c.Request.RequestURI
		// Status code
		statusCode := c.Writer.Status()
		// Client IP
		clientIP := c.ClientIP()

		// Log structure
		log.WithFields(logrus.Fields{
			"status_code":  statusCode,
			"latency_time": latencyTime,
			"client_ip":    clientIP,
			"req_method":   reqMethod,
			"req_uri":      reqURI,
		}).Info("HTTP Request")
	}
}

// CORS sets permissive cross-origin headers
func CORS() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

// RequestID attaches a unique identifier to each request
func RequestID() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Derive a simple unique ID; consider UUIDs for production use
		requestID := time.Now().UnixNano()
		c.Set("RequestID", requestID)
		c.Writer.Header().Set("X-Request-ID", time.Now().String())
		c.Next()
	}
}
