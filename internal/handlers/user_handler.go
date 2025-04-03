package handlers

import (
	"encoding/json"
	"gopark/internal/models"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// UserHandler 结构体，包含依赖项（例如日志记录器）
type UserHandler struct {
	Log *logrus.Logger
}

// NewUserHandler 创建一个新的 UserHandler 实例
func NewUserHandler(log *logrus.Logger) *UserHandler {
	return &UserHandler{Log: log}
}

// GetUser 处理获取用户信息的请求
func (h *UserHandler) GetUser(c *gin.Context) {
	user := models.User{ // Use models.User
		ID:   1,
		Name: "test",
		Mail: "test@gmail.com",
	}
	data, err := json.Marshal(user)
	if err != nil {
		h.Log.Error(err) // 使用注入的 logger
		c.JSON(500, gin.H{
			"error": "internal server error",
		})
		return
	}
	h.Log.Info(string(data)) // 使用注入的 logger
	c.JSON(200, user)
}
