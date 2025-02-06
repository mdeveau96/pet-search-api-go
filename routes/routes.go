package routes

import (
	"pet-search-backend-go/middleware"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(server *gin.Engine) {
	// Posts
	postFeed := server.Group("/feed/posts").Use(middleware.Authenticate)
	{
		postFeed.GET("/", getPosts)
		postFeed.POST("/", createPost)
		postFeed.GET("/:postId", getPost)
		postFeed.PATCH("/:postId", editPost)
		postFeed.DELETE("/:postId", deletePost)
		postFeed.POST("/:postId/like", likePost)
		postFeed.POST("/:postId/comment", postComment)
		postFeed.PATCH("/:postId/comment/:commentId", editComment)
		postFeed.DELETE("/:postId/comment/:commentId", deleteComment)
		postFeed.POST("/:postId/comment/:commentId/like", likeComment)
		postFeed.POST("/:postId/comment/:commentId/reply", postReply)
		postFeed.PATCH("/:postId/comment/:commentId/reply/:replyId", editReply)
		postFeed.DELETE("/:postId/comment/:commentId/reply/:replyId", deleteReply)
		postFeed.POST("/:postId/comment/:commentId/reply/:replyId/like", likeReply)
	}

	// Auth
	auth := server.Group("/auth")
	{
		auth.POST("/signup", signup)
		auth.POST("/login", login)
	}

	// User
	user := server.Group("/users").Use(middleware.Authenticate)
	{
		user.GET("/", getUsers)
		user.GET("/:userId", getUser)
	}
}
