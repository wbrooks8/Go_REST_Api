package middlewares

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"example.com/REST-api/utils"
	"github.com/gin-gonic/gin"
)

func TestAuthenticateRejectsMissingToken(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.Use(Authenticate)
	router.GET("/protected", func(context *gin.Context) {
		context.Status(http.StatusOK)
	})

	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodGet, "/protected", nil)
	router.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusUnauthorized {
		t.Fatalf("expected status %d, got %d", http.StatusUnauthorized, recorder.Code)
	}
}

func TestAuthenticateStoresUserIDForValidToken(t *testing.T) {
	token, err := utils.GenerateToken("user@example.com", 42)
	if err != nil {
		t.Fatalf("generate token: %v", err)
	}

	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.Use(Authenticate)
	router.GET("/protected", func(context *gin.Context) {
		context.JSON(http.StatusOK, gin.H{"userId": context.GetInt64("userId")})
	})

	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodGet, "/protected", nil)
	request.Header.Set("Authorization", token)
	router.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, recorder.Code)
	}
	if recorder.Body.String() != `{"userId":42}` {
		t.Fatalf("unexpected response: %s", recorder.Body.String())
	}
}
