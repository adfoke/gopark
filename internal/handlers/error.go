package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// ErrorResponse 定义统一的错误响应格式
type ErrorResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// RespondWithError 发送统一格式的错误响应
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

// BadRequest 处理400错误
func BadRequest(c *gin.Context, message string, log *logrus.Logger) {
	RespondWithError(c, http.StatusBadRequest, message, log)
}

// NotFound 处理404错误
func NotFound(c *gin.Context, message string, log *logrus.Logger) {
	RespondWithError(c, http.StatusNotFound, message, log)
}

// InternalServerError 处理500错误
func InternalServerError(c *gin.Context, message string, log *logrus.Logger) {
	RespondWithError(c, http.StatusInternalServerError, message, log)
}

// Unauthorized 处理401错误
func Unauthorized(c *gin.Context, message string, log *logrus.Logger) {
	RespondWithError(c, http.StatusUnauthorized, message, log)
}

// Forbidden 处理403错误
func Forbidden(c *gin.Context, message string, log *logrus.Logger) {
	RespondWithError(c, http.StatusForbidden, message, log)
}
