package routes

import (
	"gopark/internal/db"
	"gopark/internal/handlers"
	"gopark/internal/middleware"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// SetupRoutes configures and registers all application routes
func SetupRoutes(r *gin.Engine, log *logrus.Logger, dbConn *db.DB) {
	// Register global middleware
	r.Use(middleware.RequestID())
	r.Use(middleware.Logger(log))
	r.Use(middleware.CORS())

	// Create handler instances
	userHandler := handlers.NewUserHandler(log, dbConn)

	// Health check route without API versioning
	r.GET("/health", handlers.HealthCheckHandler)

	// API v1 route group
	v1 := r.Group("/api/v1")
	{
		// User management routes
		users := v1.Group("/users")
		{
			users.GET("", userHandler.GetUser)            // Query user - /api/v1/users?id=1
			users.POST("", userHandler.CreateUser)        // Create user - /api/v1/users
			users.PUT("/:id", userHandler.UpdateUser)     // Update user - /api/v1/users/1
			users.DELETE("/:id", userHandler.DeleteUser)  // Delete user - /api/v1/users/1
			users.GET("/search", userHandler.SearchUsers) // Search users - /api/v1/users/search?name=pattern
			users.GET("/list", userHandler.ListUsers)     // List users - /api/v1/users/list?limit=10&offset=0
		}
	}

	// Legacy routes retained for backward compatibility
	// Consider deprecating or removing them when appropriate
	r.GET("/user", userHandler.GetUser)
	r.POST("/user", userHandler.CreateUser)
	r.PUT("/user/:id", userHandler.UpdateUser)
	r.DELETE("/user/:id", userHandler.DeleteUser)
}
