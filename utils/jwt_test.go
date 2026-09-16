package utils

import "testing"

func TestGenerateAndVerifyToken(t *testing.T) {
	token, err := GenerateToken("user@example.com", 42)
	if err != nil {
		t.Fatalf("generate token: %v", err)
	}

	userID, err := VerifyToken(token)
	if err != nil {
		t.Fatalf("verify token: %v", err)
	}
	if userID != 42 {
		t.Fatalf("expected user ID 42, got %d", userID)
	}
}

func TestVerifyTokenRejectsMalformedToken(t *testing.T) {
	if _, err := VerifyToken("not-a-jwt"); err == nil {
		t.Fatal("expected malformed token to be rejected")
	}
}
