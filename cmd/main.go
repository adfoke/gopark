package main

import (
	"fmt"
	"gopark/config"
	"gopark/internal/routes" // Import routes package
	"gopark/internal/server" // Import server package (will be created next)
	"os"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func main() {
	// 1. 初始化日志
	log := logrus.New()
	log.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
	})
	log.SetOutput(os.Stdout) // 输出到标准输出
	// log.SetLevel(logrus.InfoLevel) // 默认 InfoLevel

	// 2. 加载配置
	cfg, err := config.LoadConfig("config") // Load from ./config directory
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// 根据配置设置日志级别
	if cfg.Debug {
		log.SetLevel(logrus.DebugLevel)
		gin.SetMode(gin.DebugMode)
		log.Info("Debug mode enabled")
	} else {
		log.SetLevel(logrus.InfoLevel)
		gin.SetMode(gin.ReleaseMode)
	}

	log.Infof("Configuration loaded: AppName=%s, Port=%d", cfg.AppName, cfg.Port)

	// 3. 创建 Gin 引擎
	r := gin.New() // Use gin.New() for more control over middleware
	r.Use(gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
		// 自定义日志格式
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
	r.Use(gin.Recovery()) // 添加 Recovery 中间件

	// 4. 注册路由
	routes.SetupRoutes(r, log) // Pass logger to routes

	// 5. 创建并启动服务器
	srv := server.NewServer(r, cfg.Port, log)
	log.Infof("Starting server on port %d", cfg.Port)
	if err := srv.Run(); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
