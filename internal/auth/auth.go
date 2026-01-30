package auth

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/alexedwards/argon2id"
	"github.com/google/uuid"
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

func MakeRefreshToken() (string, error) {
	b := make([]byte, 32)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}

func MakeToken(user_ID uuid.UUID, JWTSecret string, durInSeconds int) (string, error) {
	dur, _ := time.ParseDuration(fmt.Sprintf("%s%s", strconv.Itoa(durInSeconds), "s"))

	return MakeJWT(user_ID, JWTSecret, dur)
}

func GetBearerToken(headers http.Header) (string, error) {
	return getAuthHeader(headers, "Bearer ")
}

func GetAPIKey(headers http.Header) (string, error) {
	return getAuthHeader(headers, "ApiKey ")
}

func getAuthHeader(headers http.Header, key string) (string, error) {
	fullString := headers.Get("Authorization")
	if fullString == "" {
		return "", errors.New("Authorization Header missing")
	}

	if !strings.HasPrefix(fullString, key) {
		return "", errors.New("Invalid Authorization Header")
	}

	tokenString := strings.TrimPrefix(fullString, key)
	return tokenString, nil
}