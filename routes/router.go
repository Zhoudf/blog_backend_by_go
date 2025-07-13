package routes

import (
	"github.com/Zhoudf/blog_backend_by_go/handler"
	"github.com/Zhoudf/blog_backend_by_go/middleware"

	"github.com/gin-gonic/gin"
)

// SetupRouter 设置路由
func SetupRouter() *gin.Engine {
	// 创建Gin路由器，默认使用gin.ReleaseMode
	router := gin.Default()

	// 注册全局中间件
	router.Use(middleware.LoggerMiddleware())
	router.Use(middleware.RecoveryMiddleware())

	// API路由组
	api := router.Group("/api")
	{
		// 认证路由
		auth := api.Group("/auth")
		{
			auth.POST("/register", handler.Register) //http://localhost:8080/api/auth/register
			auth.POST("/login", handler.Login)       //http://localhost:8080/api/auth/login
		}

		// 文章路由
		posts := api.Group("/posts") //http://localhost:8080/posts
		// 评论路由
		comments := posts.Group("/:id/comments")
		{
			// 公开路由
			posts.GET("", handler.GetPosts)
			posts.GET("/:id", handler.GetPost)

			// 需要认证的文章路由
			authPosts := posts.Use(middleware.AuthMiddleware())
			{
				authPosts.POST("", handler.CreatePost)
				authPosts.PUT("/:id", handler.UpdatePost)
				authPosts.DELETE("/:id", handler.DeletePost)
			}

			{
				// 公开路由
				comments.GET("", handler.GetComments)

				// 需要认证的评论路由组
				authComments := comments.Group("")
				authComments.Use(middleware.AuthMiddleware())
				authComments.POST("", handler.CreateComment)
			}
		}
	}

	return router
}
