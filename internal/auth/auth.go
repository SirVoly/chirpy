package auth

import (
	"errors"
	"net/http"
	"strings"

	"github.com/alexedwards/argon2id"
)

// Hash the password using the argon2id.CreateHash
func HashPassword(password string) (string, error) {
	// TODO: Add Hash Parameters
	hash, err := argon2id.CreateHash(password, argon2id.DefaultParams)
	if err != nil {
		return "", err
	}

	return hash, nil
} 

// Use the argon2id.ComparePasswordAndHash function to compare the password that the user entered in the HTTP request with the password that is stored in the database.
func CheckPasswordHash(password, hash string) (bool, error) {
	return argon2id.ComparePasswordAndHash(password, hash)
}

func GetBearerToken(headers http.Header) (string, error) {
	fullString := headers.Get("Authorization")
	if fullString == "" {
		return "", errors.New("Authorization Header missing")
	}

	if !strings.HasPrefix(fullString, "Bearer ") {
		return "", errors.New("Invalid Authorization Header")
	}

	tokenString := strings.TrimPrefix(fullString, "Bearer ")
	return tokenString, nil
}