package middleware

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

const secretKey = "pet-search-api-secret"

func Authenticate(context *gin.Context) {
	authHeader := context.Request.Header.Get("Authorization")
	if authHeader == "" {
		context.JSON(http.StatusUnauthorized, gin.H{"message": "Invalid email or password", "error": "authheader"})
		context.Abort()
	}
	token := strings.Split(authHeader, " ")[1]
	if token == "" {
		context.JSON(http.StatusUnauthorized, gin.H{"message": "Invalid email or password", "error": "token"})
		context.Abort()
	}
	decodedToken, err := jwt.Parse(token, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", t.Header["alg"])
		}
		return []byte(secretKey), nil
	})
	if err != nil {
		context.JSON(http.StatusUnauthorized, gin.H{"message": "Invalid email or password", "error": "decoded"})
		context.Abort()
	}
	if claims, ok := decodedToken.Claims.(jwt.MapClaims); ok {
		context.Request.Header.Set("userId", claims["sub"].(string))
		context.Next()
	} else {
		context.JSON(http.StatusUnauthorized, gin.H{"message": "Invalid email or password", "error": "claims"})
		context.Abort()
	}
}
