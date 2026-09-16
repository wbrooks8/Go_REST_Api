package middlewares

import (
	"net/http"

	"example.com/REST-api/utils"
	"github.com/gin-gonic/gin"
)

// Authenticate is middleware: it runs before protected handlers, verifies the
// bearer token, and stores the authenticated user's ID in Gin's context.
func Authenticate(context *gin.Context) {
	token := context.Request.Header.Get("Authorization")

	if token == "" {
		context.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "Not authorized."})
		return
	}

	userId, err := utils.VerifyToken(token)

	if err != nil {
		context.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "Not authorized"})
		return
	}

	context.Set("userId", userId)

	context.Next()
}
