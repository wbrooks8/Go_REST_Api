package utils

import "golang.org/x/crypto/bcrypt"

// HashPassword uses bcrypt so the original password is never stored.
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)

	return string(bytes), err
}

// CheckPasswordHash returns false for an incorrect or malformed hash.
func CheckPasswordHash(password, hashedPassword string) bool {

	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	return err == nil
}
