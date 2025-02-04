package routes

import (
	"net/http"
	"pet-search-backend-go/models"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

func createToken(userId primitive.ObjectID) (string, error) {
	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub":    userId,
		"iss":    "pet-search",
		"exp":    time.Now().Add(time.Hour).Unix(),
		"iat":    time.Now().Unix(),
	})
	signingKey := []byte("pet-search-api-secret")
	tokenString, err := claims.SignedString(signingKey)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func verifyPassword(hashedPassword, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	return err == nil
}

func signup(context *gin.Context) {
	var newUser models.User
	err := context.ShouldBindJSON(&newUser)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Could not parse request data"})
		return
	}
	createdUser, err := newUser.AddUser()
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Could not create post"})
		return
	}
	context.JSON(http.StatusCreated, gin.H{"message": "User created", "user": createdUser})
}

func login(context *gin.Context) {
	var credentials models.Login
	err := context.ShouldBindJSON(&credentials)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Could not parse request data"})
		return
	}
	user, err := models.FindUserByEmail(credentials.Email)
	if err != nil {
		context.JSON(http.StatusNotFound, gin.H{"message": "Could not find user"})
		return
	}
	if credentials.Email != user.Email {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Invalid Username or Password"})
		return
	}
	if !(verifyPassword(user.Password, credentials.Password)) {
		context.JSON(http.StatusBadRequest, gin.H{"message": "Invalid Username or Password"})
		return
	}
	token, err := createToken(user.ID)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": "Could not log in user", "error": err})
		return
	}
	context.JSON(http.StatusAccepted, gin.H{"message": "Login Successful", "token": token, "user": user.Email})
}
