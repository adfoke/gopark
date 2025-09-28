package server

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// Server encapsulates the HTTP server and its dependencies
type Server struct {
	httpServer *http.Server
	log        *logrus.Logger
}

// NewServer creates a new Server instance
func NewServer(router *gin.Engine, port int, log *logrus.Logger) *Server {
	addr := fmt.Sprintf(":%d", port)

	httpServer := &http.Server{
		Addr:           addr,
		Handler:        router,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		IdleTimeout:    120 * time.Second,
		MaxHeaderBytes: 1 << 20, // 1 MB
	}

	return &Server{
		httpServer: httpServer,
		log:        log,
	}
}

// Run starts the HTTP server and supports graceful shutdown
func (s *Server) Run() error {
	// Start the server in a dedicated goroutine
	go func() {
		s.log.Infof("HTTP server listening on %s", s.httpServer.Addr)
		if err := s.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			s.log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// Wait for interrupt signals to shut down gracefully
	quit := make(chan os.Signal, 1)
	// kill (no args) sends syscall.SIGTERM
	// kill -2 sends syscall.SIGINT
	// kill -9 sends syscall.SIGKILL, which cannot be caught
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	s.log.Info("Shutting down server...")

	// Allow up to 5 seconds to complete shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := s.Shutdown(ctx); err != nil {
		s.log.Fatalf("Server forced to shutdown: %v", err)
	}

	s.log.Info("Server exiting")
	return nil
}

// Shutdown stops the server gracefully
func (s *Server) Shutdown(ctx context.Context) error {
	s.log.Info("Shutting down server...")
	return s.httpServer.Shutdown(ctx)
}
