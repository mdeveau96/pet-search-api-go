package routes

import (
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(server *gin.Engine) {
	// Posts
	server.GET("/feed/posts", getPosts)
	server.GET("/feed/posts/:id", getPost)
	server.POST("/feed/posts", createPost)
	server.PATCH("/feed/posts/:id", updatePost)
	server.DELETE("/feed/posts/:id", deletePost)

	server.POST("/feed/posts/:id/like", likePost)

	// Auth
	server.POST("/auth/signup", signup)
	server.POST("/auth/login", login)
}
