package routes

import (
	"gopark/internal/db"
	"gopark/internal/handlers"
	"gopark/internal/middleware"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// SetupRoutes 配置和注册所有应用路由
func SetupRoutes(r *gin.Engine, log *logrus.Logger, dbConn *db.DB) {
	// 添加全局中间件
	r.Use(middleware.RequestID())
	r.Use(middleware.Logger(log))
	r.Use(middleware.CORS())

	// 创建 handler 实例
	userHandler := handlers.NewUserHandler(log, dbConn)
	helloHandler := handlers.NewHelloHandler(log)

	// 健康检查路由 - 不需要API版本
	r.GET("/health", handlers.HealthCheckHandler)

	// API v1 路由组
	v1 := r.Group("/api/v1")
	{
		// Hello 路由
		v1.GET("/hello", helloHandler.SayHello)

		// 用户管理路由
		users := v1.Group("/users")
		{
			users.GET("", userHandler.GetUser)            // 查询用户 - /api/v1/users?id=1
			users.POST("", userHandler.CreateUser)        // 创建用户 - /api/v1/users
			users.PUT("/:id", userHandler.UpdateUser)     // 更新用户 - /api/v1/users/1
			users.DELETE("/:id", userHandler.DeleteUser)  // 删除用户 - /api/v1/users/1
			users.GET("/search", userHandler.SearchUsers) // 搜索用户 - /api/v1/users/search?name=pattern
			users.GET("/list", userHandler.ListUsers)     // 列出用户 - /api/v1/users/list?limit=10&offset=0
		}
	}

	// 为了向后兼容，保留旧路由
	// 在实际项目中，可以考虑添加弃用警告或在适当时候移除
	r.GET("/user", userHandler.GetUser)
	r.POST("/user", userHandler.CreateUser)
	r.PUT("/user/:id", userHandler.UpdateUser)
	r.DELETE("/user/:id", userHandler.DeleteUser)
	r.GET("/hello", helloHandler.SayHello)
}
