package auth

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

func MakeJWT(userID uuid.UUID, tokenSecret string, expiresIn time.Duration) (string, error) {

	// Create the Claims
	claims := &jwt.RegisteredClaims{
		Issuer:    "chirpy",
		IssuedAt: jwt.NewNumericDate(time.Now()),
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(expiresIn)),
		Subject: userID.String(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	ss, err := token.SignedString([]byte(tokenSecret))

	return ss, err
}

func ValidateJWT(tokenString, tokenSecret string) (uuid.UUID, error) {

	token, err := jwt.ParseWithClaims(tokenString, &jwt.RegisteredClaims{}, func(token *jwt.Token) (any, error) {
		return []byte(tokenSecret), nil
	})

	if err != nil {
		return uuid.UUID{}, err
	}

	expiresAt, err := token.Claims.GetExpirationTime()
	
	if  err != nil || expiresAt.Time.Unix() < time.Now().Unix() {
		return uuid.UUID{}, fmt.Errorf("Token expired!")
	}

	userID, err := token.Claims.GetSubject()

	if err != nil {
		return uuid.UUID{}, err
	}

	return uuid.MustParse(userID), nil
}