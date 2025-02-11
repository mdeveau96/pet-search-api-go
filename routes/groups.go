package routes

import (
	"net/http"
	"pet-search-backend-go/models"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func getGroups(context *gin.Context) {
	users, err := models.FindAllUsers()
	if err != nil {
		context.JSON(http.StatusNotFound, gin.H{"message": "Could not find users", "error": err})
		return
	}
	context.JSON(http.StatusOK, gin.H{"users": users})
}

func getGroup(context *gin.Context) {
	userId, _ := primitive.ObjectIDFromHex(context.Request.Header.Get("userId"))
	filter := bson.D{{Key: "_id", Value: userId}}
	user, err := models.FindUser(filter)
	if err != nil {
		context.JSON(http.StatusNotFound, gin.H{"message": "Could not find user"})
		return
	}
	context.JSON(http.StatusOK, gin.H{"user": user})
}
