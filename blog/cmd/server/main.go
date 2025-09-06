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
	"fmt"
	"github.com/spf13/viper"
)

func main() {
	// 1. 加载配置
	config.InitConfig("setting.yaml")

	logger.InitLogger(false)

	// 2. 初始化数据库
	err := db.InitDB()
	if err != nil {
		panic(fmt.Errorf("InitDB error: %w", err))
	}

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

	// 5. 设置路由并启动服务器
	// 确保 NewRouter 函数现在接收 commentHandler
	r := api.NewRouter(userHandler, postHandler, commentHandler)
	if err := r.Run(viper.GetString("server.port")); err != nil {
		panic(fmt.Errorf("failed to start server: %w", err))
	}
}
