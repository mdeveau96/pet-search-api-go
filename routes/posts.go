package routes

import (
	"net/http"
	"pet-search-backend-go/models"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
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
	userId, _ := primitive.ObjectIDFromHex(context.Request.Header.Get("userId"))
	post.Creator = userId
	newPost, err := post.Create()
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Could not create post"})
		return
	}
	filter := bson.D{{Key: "_id", Value: userId}}
	user, err := models.FindUser(filter)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": "Could not attach post to user account"})
		return
	}
	user.AddPost(newPost)
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

func deletePost(context *gin.Context) {
	postId, err := primitive.ObjectIDFromHex(context.Param("id"))
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Could not find post"})
		return
	}
	result, err := models.Delete(postId)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": "Unable to delete post", "error": err})
		return
	}
	context.JSON(http.StatusOK, gin.H{"message": "Post deleted", "post": result})
}

func likePost(context *gin.Context) {
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
	userId, _ := primitive.ObjectIDFromHex(context.Request.Header.Get("userId"))
	err = post.Like(userId)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": "Unable to like post"})
	}
	context.JSON(http.StatusOK, gin.H{"message": "Post liked", "post": post})
}

func postComment(context *gin.Context) {
	var newComment models.Comment
	err := context.ShouldBindJSON(&newComment)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Could not parse request data"})
		return
	}
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
	userId, _ := primitive.ObjectIDFromHex(context.Request.Header.Get("userId"))
	newComment.Creator = userId
	err = post.AddComment(newComment)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": "Unable to add coment", "error": err})
		return
	}
	context.JSON(http.StatusOK, gin.H{"message": "Comment added", "post": post})
}
