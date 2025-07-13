package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Zhoudf/blog_backend_by_go/config"
	"github.com/Zhoudf/blog_backend_by_go/routes"
)

func main() {
	// 初始化数据库
	if err := config.InitDB(); err != nil {
		log.Fatalf("数据库初始化失败: %v", err)
	}

	// 设置路由
	router := routes.SetupRouter()

	// 从环境变量获取端口，默认8080
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// 创建HTTP服务器
	srv := &http.Server{
		Addr:    ":" + port,
		Handler: router,
	}

	// 启动服务器（非阻塞）
	go func() {
		log.Printf("服务器正在端口 %s 上运行...", port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("服务器启动失败: %v", err)
		}
	}()

	// 等待中断信号以优雅关闭服务器
	quit := make(chan os.Signal, 1)
	// 接收syscall.SIGINT和syscall.SIGTERM信号
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("正在关闭服务器...")

	// 设置5秒超时关闭服务器
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("服务器关闭失败:", err)
	}

	log.Println("服务器已关闭")
}
