package server

import (
	"fmt"
	"net/http"
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
		MaxHeaderBytes: 1 << 20, // 1 MB
	}

	return &Server{
		httpServer: httpServer,
		log:        log,
	}
}

// Run 启动 HTTP 服务器
func (s *Server) Run() error {
	s.log.Infof("HTTP server listening on %s", s.httpServer.Addr)
	return s.httpServer.ListenAndServe()
}

// 可以添加优雅关闭的逻辑
// func (s *Server) Shutdown(ctx context.Context) error {
// 	 s.log.Info("Shutting down server...")
// 	 return s.httpServer.Shutdown(ctx)
// }
