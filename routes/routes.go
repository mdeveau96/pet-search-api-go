package routes

import (
	"pet-search-backend-go/middleware"
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(server *gin.Engine) {
	// Posts
	server.GET("/feed/posts", middleware.Authenticate, getPosts)
	server.GET("/feed/posts/:id", middleware.Authenticate, getPost)
	server.POST("/feed/posts", middleware.Authenticate, createPost)
	server.PATCH("/feed/posts/:id", middleware.Authenticate, updatePost)
	server.DELETE("/feed/posts/:id", middleware.Authenticate, deletePost)

	server.POST("/feed/posts/:id/like", middleware.Authenticate, likePost)

	// Auth
	server.POST("/auth/signup", signup)
	server.POST("/auth/login", login)
}
