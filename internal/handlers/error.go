package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// ErrorResponse defines the common error payload
type ErrorResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// RespondWithError sends a consistent error response
func RespondWithError(c *gin.Context, statusCode int, message string, log *logrus.Logger) {
	log.WithFields(logrus.Fields{
		"status_code": statusCode,
		"error":       message,
		"path":        c.Request.URL.Path,
		"method":      c.Request.Method,
	}).Error("Request error")

	c.JSON(statusCode, ErrorResponse{
		Code:    statusCode,
		Message: message,
	})
}

// BadRequest handles a 400 Bad Request response
func BadRequest(c *gin.Context, message string, log *logrus.Logger) {
	RespondWithError(c, http.StatusBadRequest, message, log)
}

// NotFound handles a 404 Not Found response
func NotFound(c *gin.Context, message string, log *logrus.Logger) {
	RespondWithError(c, http.StatusNotFound, message, log)
}

// InternalServerError handles a 500 Internal Server Error response
func InternalServerError(c *gin.Context, message string, log *logrus.Logger) {
	RespondWithError(c, http.StatusInternalServerError, message, log)
}

// Unauthorized handles a 401 Unauthorized response
func Unauthorized(c *gin.Context, message string, log *logrus.Logger) {
	RespondWithError(c, http.StatusUnauthorized, message, log)
}

// Forbidden handles a 403 Forbidden response
func Forbidden(c *gin.Context, message string, log *logrus.Logger) {
	RespondWithError(c, http.StatusForbidden, message, log)
}
