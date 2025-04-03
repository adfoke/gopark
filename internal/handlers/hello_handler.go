package handlers

import (
	"gopark/internal/hello"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// HelloHandler 结构体，包含依赖项
type HelloHandler struct {
	Log *logrus.Logger
}

// NewHelloHandler 创建一个新的 HelloHandler 实例
func NewHelloHandler(log *logrus.Logger) *HelloHandler {
	return &HelloHandler{Log: log}
}

// SayHello 处理 /hello 请求
func (h *HelloHandler) SayHello(c *gin.Context) {
	message := hello.SayHello()
	h.Log.Info(message) // 使用注入的 logger
	c.JSON(200, gin.H{
		"message": message,
	})
}
