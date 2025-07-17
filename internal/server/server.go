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

// Server 结构体封装了 HTTP 服务器及其依赖
type Server struct {
	httpServer *http.Server
	log        *logrus.Logger
}

// NewServer 创建一个新的 Server 实例
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

// Run 启动 HTTP 服务器并支持优雅关闭
func (s *Server) Run() error {
	// 在单独的 goroutine 中启动服务器
	go func() {
		s.log.Infof("HTTP server listening on %s", s.httpServer.Addr)
		if err := s.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			s.log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// 等待中断信号以优雅地关闭服务器
	quit := make(chan os.Signal, 1)
	// kill (无参数) 默认发送 syscall.SIGTERM
	// kill -2 发送 syscall.SIGINT
	// kill -9 发送 syscall.SIGKILL，但无法被捕获，所以不需要添加
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	s.log.Info("Shutting down server...")

	// 设置 5 秒的超时时间来关闭服务器
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := s.Shutdown(ctx); err != nil {
		s.log.Fatalf("Server forced to shutdown: %v", err)
	}

	s.log.Info("Server exiting")
	return nil
}

// Shutdown 优雅地关闭服务器
func (s *Server) Shutdown(ctx context.Context) error {
	s.log.Info("Shutting down server...")
	return s.httpServer.Shutdown(ctx)
}
