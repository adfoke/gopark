package handlers

import (
	"github.com/gin-gonic/gin"
)

// HealthCheckHandler handles the health check request
func HealthCheckHandler(c *gin.Context) {
	c.JSON(200, gin.H{
		"status": "ok",
	})
}
