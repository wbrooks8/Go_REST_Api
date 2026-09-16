package utils

import "testing"

func TestHashPasswordAndCheckPasswordHash(t *testing.T) {
	hashedPassword, err := HashPassword("correct horse battery staple")
	if err != nil {
		t.Fatalf("hash password: %v", err)
	}
	if hashedPassword == "correct horse battery staple" {
		t.Fatal("password must not be stored as plain text")
	}
	if !CheckPasswordHash("correct horse battery staple", hashedPassword) {
		t.Fatal("expected password hash to validate")
	}
	if CheckPasswordHash("wrong password", hashedPassword) {
		t.Fatal("expected wrong password to fail")
	}
}
