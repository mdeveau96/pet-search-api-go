package routes

import (
	"net/http"
	"pet-search-backend-go/models"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type params struct {
	PostId    primitive.ObjectID
	CommentId primitive.ObjectID
	replyId   primitive.ObjectID
}

func getIdsFromParams(context *gin.Context) (params, error) {
	postId, err := primitive.ObjectIDFromHex(context.Param("postId"))
	if err != nil {
		// context.JSON(http.StatusBadRequest, gin.H{"message": "Could not parse postId from request data"})
		return params{PostId: primitive.NilObjectID, CommentId: primitive.NilObjectID, replyId: primitive.NilObjectID}, err
	}
	commentId, err := primitive.ObjectIDFromHex(context.Param("commentId"))
	if err != nil {
		// context.JSON(http.StatusBadRequest, gin.H{"message": "Could not parse commentId from request data"})
		return params{PostId: postId, CommentId: primitive.NilObjectID, replyId: primitive.NilObjectID}, nil
	}
	replyId, err := primitive.ObjectIDFromHex(context.Param("replyId"))
	if err != nil {
		// context.JSON(http.StatusBadRequest, gin.H{"message": "Could not parse replyId from request data"})
		return params{PostId: postId, CommentId: commentId, replyId: primitive.NilObjectID}, nil
	}
	return params{PostId: postId, CommentId: commentId, replyId: replyId}, nil
}

func updateUserPosts(context *gin.Context, post models.Post) {
	userId, _ := primitive.ObjectIDFromHex(context.Request.Header.Get("userId"))
	filter := bson.D{{Key: "_id", Value: userId}}
	user, err := models.FindUser(filter)
	if err != nil {
		context.JSON(http.StatusNotFound, gin.H{"message": "Could not find user"})
		return
	}
	err = user.UpdateUserPosts(post)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": "Could not attach post to user account"})
		return
	}
}

func deleteUserPost(context *gin.Context, postId primitive.ObjectID) {
	userId, _ := primitive.ObjectIDFromHex(context.Request.Header.Get("userId"))
	filter := bson.D{{Key: "_id", Value: userId}}
	user, err := models.FindUser(filter)
	if err != nil {
		context.JSON(http.StatusNotFound, gin.H{"message": "Could not find user"})
		return
	}
	err = user.DeleteUserPost(postId)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": "Could not remove post from user account"})
		return
	}
}

func getPosts(context *gin.Context) {
	posts, err := models.FindAllPosts()
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": "Could not fetch posts. Try again later", "error": err})
		return
	}
	context.JSON(http.StatusOK, posts)
}

func getPost(context *gin.Context) {
	params, err := getIdsFromParams(context)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Could not parse request data"})
		return
	}
	post, err := models.FindPost(params.PostId)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": "Could not fetch post"})
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

func editPost(context *gin.Context) {
	var updatedPost models.Post
	err := context.ShouldBindJSON(&updatedPost)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Could not parse request data"})
		return
	}
	params, err := getIdsFromParams(context)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Could not parse request data"})
		return
	}
	post, err := models.FindPost(params.PostId)
	if err != nil {
		context.JSON(http.StatusNotFound, gin.H{"message": "Could not find post"})
		return
	}
	updatedPost.CreatedAt = post.CreatedAt
	updatedPost.ID = params.PostId
	result, err := updatedPost.Update(updatedPost)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": "Could not update post"})
		return
	}
	updateUserPosts(context, result)
	context.JSON(http.StatusOK, gin.H{"message": "Post Updated", "post": result})
}

func deletePost(context *gin.Context) {
	params, err := getIdsFromParams(context)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Could not parse request data"})
		return
	}
	result, err := models.Delete(params.PostId)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": "Unable to delete post"})
		return
	}
	deleteUserPost(context, params.PostId)
	context.JSON(http.StatusOK, gin.H{"message": "Post deleted", "post": result})
}

func likePost(context *gin.Context) {
	params, err := getIdsFromParams(context)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Could not parse request data"})
		return
	}
	post, err := models.FindPost(params.PostId)
	if err != nil {
		context.JSON(http.StatusNotFound, gin.H{"message": "Could not fetch post"})
		return
	}
	userId, _ := primitive.ObjectIDFromHex(context.Request.Header.Get("userId"))
	result, err := post.Like(userId)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": "Unable to like post"})
	}
	updateUserPosts(context, result)
	context.JSON(http.StatusOK, gin.H{"message": "Post liked", "post": result})
}

func postComment(context *gin.Context) {
	var newComment models.Comment
	err := context.ShouldBindJSON(&newComment)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Could not parse request data"})
		return
	}
	params, err := getIdsFromParams(context)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Could not parse request data"})
		return
	}
	post, err := models.FindPost(params.PostId)
	if err != nil {
		context.JSON(http.StatusNotFound, gin.H{"message": "Could not fetch post"})
		return
	}
	userId, _ := primitive.ObjectIDFromHex(context.Request.Header.Get("userId"))
	newComment.Creator = userId
	result, err := post.AddComment(newComment)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": "Unable to add comment", "error": err})
		return
	}
	updateUserPosts(context, result)
	context.JSON(http.StatusCreated, gin.H{"message": "Comment added", "post": result})
}

func likeComment(context *gin.Context) {
	params, err := getIdsFromParams(context)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Could not parse request data"})
		return
	}
	userId, err := primitive.ObjectIDFromHex(context.Request.Header.Get("userId"))
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Could not read user id header"})
		return
	}
	post, err := models.FindPost(params.PostId)
	if err != nil {
		context.JSON(http.StatusNotFound, gin.H{"message": "Could not fetch post"})
		return
	}
	result, err := post.LikeComment(params.CommentId, userId)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": "Could not fetch post"})
		return
	}
	updateUserPosts(context, result)
	context.JSON(http.StatusOK, gin.H{"message": "Comment liked", "post": result})
}

func editComment(context *gin.Context) {
	var updatedComment models.Comment
	err := context.ShouldBindJSON(&updatedComment)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Could not parse request data", "error": err})
		return
	}
	params, err := getIdsFromParams(context)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Could not parse request data"})
		return
	}
	post, err := models.FindPost(params.PostId)
	if err != nil {
		context.JSON(http.StatusNotFound, gin.H{"message": "Could not fetch post"})
		return
	}
	updatedComment.ID = params.CommentId
	result, err := post.UpdateComment(updatedComment)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": "Could not update comment"})
		return
	}
	updateUserPosts(context, result)
	context.JSON(http.StatusOK, gin.H{"message": "Updated comment", "post": result})
}

func postReply(context *gin.Context) {
	var reply models.Reply
	err := context.ShouldBindJSON(&reply)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Could not parse request data", "error": err})
		return
	}
	params, err := getIdsFromParams(context)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Could not parse request data"})
		return
	}
	userId, err := primitive.ObjectIDFromHex(context.Request.Header.Get("userId"))
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Could not read user id header"})
		return
	}
	post, err := models.FindPost(params.PostId)
	if err != nil {
		context.JSON(http.StatusNotFound, gin.H{"message": "Could not fetch post"})
		return
	}
	reply.Creator = userId
	result, err := post.ReplyToComment(params.CommentId, reply)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": "Could not reply to comment"})
		return
	}
	updateUserPosts(context, result)
	context.JSON(http.StatusOK, gin.H{"message": "Reply posted", "post": result})
}

func editReply(context *gin.Context) {
	var updatedReply models.Reply
	err := context.ShouldBindJSON(&updatedReply)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Could not parse request data", "error": err})
		return
	}
	params, err := getIdsFromParams(context)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Could not parse request data"})
		return
	}
	post, err := models.FindPost(params.PostId)
	if err != nil {
		context.JSON(http.StatusNotFound, gin.H{"message": "Could not fetch post"})
		return
	}
	updatedReply.ID = params.replyId
	result, err := post.EditReply(params.CommentId, updatedReply)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": "Could not edit reply"})
		return
	}
	updateUserPosts(context, result)
	context.JSON(http.StatusOK, gin.H{"message": "Editted reply", "post": result})
}

func likeReply(context *gin.Context) {
	params, err := getIdsFromParams(context)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Could not parse request data"})
		return
	}
	userId, err := primitive.ObjectIDFromHex(context.Request.Header.Get("userId"))
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Could not read user id header"})
		return
	}
	post, err := models.FindPost(params.PostId)
	if err != nil {
		context.JSON(http.StatusNotFound, gin.H{"message": "Could not fetch post"})
		return
	}
	result, err := post.LikeReply(params.CommentId, userId, params.replyId)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": "Could not like reply"})
		return
	}
	updateUserPosts(context, result)
	context.JSON(http.StatusOK, gin.H{"message": "Liked reply", "post": result})
}

func deleteReply(context *gin.Context) {
	params, err := getIdsFromParams(context)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Could not parse request data"})
		return
	}
	post, err := models.FindPost(params.PostId)
	if err != nil {
		context.JSON(http.StatusNotFound, gin.H{"message": "Could not fetch post"})
		return
	}
	result, err := post.DeleteReply(params.CommentId, params.replyId)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": "Could not delete reply"})
		return
	}
	updateUserPosts(context, result)
	context.JSON(http.StatusOK, gin.H{"message": "Deleted reply", "post": result})
}

func deleteComment(context *gin.Context) {
	params, err := getIdsFromParams(context)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Could not parse request data"})
		return
	}
	post, err := models.FindPost(params.PostId)
	if err != nil {
		context.JSON(http.StatusNotFound, gin.H{"message": "Could not fetch post"})
		return
	}
	result, err := post.DeleteComment(params.CommentId)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": "Could not delete comment"})
		return
	}
	updateUserPosts(context, result)
	context.JSON(http.StatusOK, gin.H{"message": "Deleted comment", "post": result})
}
