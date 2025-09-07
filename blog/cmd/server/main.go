package main

import (
	"blog/internal/api"
	"blog/internal/api/handler"
	"blog/internal/pkg/config"
	"blog/internal/pkg/db"
	"blog/internal/pkg/jwt"
	"blog/internal/pkg/logger"
	"blog/internal/repository"
	"blog/internal/service"
	"context"
	"fmt"
	"github.com/spf13/viper"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	// 1. 加载配置
	config.InitConfig("/Users/owen/workspace/web3/web3-tasks/blog/setting.yaml")

	logger.InitLogger(false)

	// 2. 初始化数据库
	db.InitDB()

	// 获取primary db
	primaryDB := db.Dbs["primary"]

	// 3. 初始化JWT
	jwt.Setup(viper.GetString("jwt.secret"))

	// 4. 依赖注入：从Repository层 -> Service层 -> Handler层
	userRepo := repository.NewUserRepository(primaryDB)
	postRepo := repository.NewPostRepository(primaryDB)
	commentRepo := repository.NewCommentRepository(primaryDB) // 新增

	userService := service.NewUserService(
		userRepo,
		viper.GetString("jwt.issuer"),
		viper.GetInt("jwt.expire_hours"),
	)
	postService := service.NewPostService(postRepo)
	// 注意：CommentService 依赖 CommentRepository 和 PostRepository
	commentService := service.NewCommentService(commentRepo, postRepo) // 新增

	userHandler := handler.NewUserHandler(userService)
	postHandler := handler.NewPostHandler(postService)
	commentHandler := handler.NewCommentHandler(commentService) // 新增

	// 5. 设置路由
	r := api.NewRouter(userHandler, postHandler, commentHandler)

	// 6. 启动服务并实现优雅停机
	startServer(r)
}

func startServer(router http.Handler) {
	// 创建 HTTP 服务器实例
	srv := &http.Server{
		Addr:    viper.GetString("server.port"),
		Handler: router,
	}

	// 在一个新的 goroutine 中启动服务器，以避免阻塞主线程
	go func() {
		log.Printf("服务器正在端口 %s 上启动...\n", srv.Addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			panic(fmt.Errorf("服务器启动失败: %w", err))
		}
	}()

	// 创建一个通道来接收系统信号
	quit := make(chan os.Signal, 1)
	// 监听 SIGINT 和 SIGTERM 信号
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	// 阻塞主线程，直到接收到信号
	<-quit
	log.Println("接收到关闭信号，正在准备关闭服务器...")

	// 创建一个有超时的上下文，用于通知服务器在5秒内完成正在处理的请求
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// 调用 Shutdown() 优雅地关闭服务器
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("服务器强制关闭出错: %v", err)
	}
	log.Println("HTTP 服务器已成功关闭。")

	// 在 HTTP 服务器关闭后，关闭所有数据库连接
	db.CloseDBs()

	log.Println("服务已成功优雅退出。")
}
