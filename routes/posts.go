package routes

import (
	"net/http"
	"pet-search-backend-go/models"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func getPosts(context *gin.Context) {
	posts, err := models.FindAllPosts()
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": "Could not fetch posts. Try again later", "error": err})
		return
	}
	context.JSON(http.StatusOK, posts)
}

func getPost(context *gin.Context) {
	postId, err := primitive.ObjectIDFromHex(context.Param("id"))
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": "Could not read post id param."})
		return
	}
	post, err := models.FindPost(postId)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": "Could not fetch post."})
		return
	}
	context.JSON(http.StatusOK, post)
}

func createPost(context *gin.Context) {
	var post models.Post
	err := context.ShouldBindJSON(&post)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Could not parse request data"})
		return
	}
	newPost, err := post.Create()
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Could not create post"})
		return
	}
	context.JSON(http.StatusCreated, gin.H{"message": "Post created", "post": newPost})
}

func updatePost(context *gin.Context) {
	var updatedPost models.Post
	err := context.ShouldBindJSON(&updatedPost)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Could not parse request data", "error": err})
		return
	}
	postId, err := primitive.ObjectIDFromHex(context.Param("id"))
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Could not find post"})
		return
	}
	post, err := models.FindPost(postId)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Could not find post"})
		return
	}
	updatedPost.CreatedAt = post.CreatedAt
	updatedPost.ID = postId
	result, err := updatedPost.Update(updatedPost)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": "Could not update post", "error": err})
		return
	}
	context.JSON(http.StatusOK, gin.H{"message": "Post Updated", "post": result})
}
