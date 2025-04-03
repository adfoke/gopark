package handlers

import (
	"github.com/gin-gonic/gin"
)

// HealthCheckHandler 处理健康检查请求
func HealthCheckHandler(c *gin.Context) {
	c.JSON(200, gin.H{
		"status": "ok",
	})
}
