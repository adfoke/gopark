package routes

import (
	"gopark/internal/db"
	"gopark/internal/handlers"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// SetupRoutes 配置和注册所有应用路由
func SetupRoutes(r *gin.Engine, log *logrus.Logger, dbConn *db.DB) {
	// 创建 handler 实例
	userHandler := handlers.NewUserHandler(log, dbConn)
	helloHandler := handlers.NewHelloHandler(log)

	// 注册路由
	r.GET("/health", handlers.HealthCheckHandler) // Health check 不需要特定 handler 实例
	r.GET("/user", userHandler.GetUser)
	r.GET("/hello", helloHandler.SayHello)
	
	// 添加用户管理路由
	r.POST("/user", userHandler.CreateUser)
	r.PUT("/user/:id", userHandler.UpdateUser)
	r.DELETE("/user/:id", userHandler.DeleteUser)
}
