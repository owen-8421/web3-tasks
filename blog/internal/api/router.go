package api

import (
	"blog/internal/api/handler"
	"blog/internal/api/middleware"
	"github.com/gin-gonic/gin"
)

func NewRouter(
	userHandler *handler.UserHandler,
	postHandler *handler.PostHandler,
	commentHandler *handler.CommentHandler, // 确保接收 commentHandler
) *gin.Engine {
	r := gin.Default()

	apiV1 := r.Group("/api/v1")
	{
		// --- 用户认证 ---
		apiV1.POST("/register", userHandler.Register)
		apiV1.POST("/login", userHandler.Login)

		// --- 文章路由 (公共) ---
		postsPublic := apiV1.Group("/posts")
		{
			postsPublic.GET("", postHandler.GetAll)
			postsPublic.GET("/:id", postHandler.GetByID)
		}

		// --- 评论路由 (公共) ---
		// 获取某篇文章的所有评论
		apiV1.GET("/posts/:id/comments", commentHandler.GetByPostID)

		// --- 需要认证的路由组 ---
		authed := apiV1.Group("/")
		authed.Use(middleware.JWTAuth())
		{
			// 文章管理 (创建、更新、删除)
			authed.POST("/posts", postHandler.Create)
			authed.PUT("/posts/:id", postHandler.Update)
			authed.DELETE("/posts/:id", postHandler.Delete)

			// 评论管理 (创建、删除)
			authed.POST("/posts/:id/comments", commentHandler.Create)
			// 注意: 删除评论的路由最好是针对评论本身的资源
			authed.DELETE("/comments/:id", commentHandler.Delete)
		}
	}

	return r
}
