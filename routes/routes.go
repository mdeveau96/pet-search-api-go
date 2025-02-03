package routes

import (
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(server *gin.Engine) {
	server.GET("/feed/posts", getPosts)
	server.GET("/feed/posts/:id", getPost)

	server.POST("/feed/posts", createPost)
	server.PATCH("/feed/posts/:id", updatePost)
}
