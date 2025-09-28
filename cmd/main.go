package main

import (
	"context"
	"fmt"
	"gopark/config"
	"gopark/internal/db"     // Import database package
	"gopark/internal/routes" // Import routes package
	"gopark/internal/server" // Import server package (will be created next)
	"os"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func main() {
	// Initialize logger
	log := logrus.New()
	log.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
	})
	log.SetOutput(os.Stdout)
	// log.SetLevel(logrus.InfoLevel)

	// Load configuration
	cfg, err := config.LoadConfig("config") // Load from ./config directory
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Configure logging based on debug flag
	if cfg.Debug {
		log.SetLevel(logrus.DebugLevel)
		gin.SetMode(gin.DebugMode)
		log.Info("Debug mode enabled")
	} else {
		log.SetLevel(logrus.InfoLevel)
		gin.SetMode(gin.ReleaseMode)
	}

	log.Infof("Configuration loaded: AppName=%s, Port=%d", cfg.AppName, cfg.Port)

	// Create Gin engine
	r := gin.New() // Use gin.New() for more control over middleware
	r.Use(gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
		// Custom log format
		return fmt.Sprintf("%s - [%s] \"%s %s %s %d %s \"%s\" %s\"\n",
			param.ClientIP,
			param.TimeStamp.Format("2006/01/02 - 15:04:05"),
			param.Method,
			param.Path,
			param.Request.Proto,
			param.StatusCode,
			param.Latency,
			param.Request.UserAgent(),
			param.ErrorMessage,
		)
	}))
	r.Use(gin.Recovery())

	// Initialize database connection
	dbConn, err := db.NewDB(cfg, log)
	if err != nil {
		log.Fatalf("Failed to initialize database connection: %v", err)
	}
	defer dbConn.Close()

	// Run database migrations
	migrationManager := db.NewMigrationManager(dbConn, log)
	if err := migrationManager.RunMigrations(context.Background(), "internal/migrations"); err != nil {
		log.Fatalf("Failed to run database migrations: %v", err)
	}
	log.Info("Database migrations completed successfully")

	// Register routes
	routes.SetupRoutes(r, log, dbConn)

	// Create and start server
	srv := server.NewServer(r, cfg.Port, log)
	log.Infof("Starting server on port %d", cfg.Port)
	if err := srv.Run(); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
