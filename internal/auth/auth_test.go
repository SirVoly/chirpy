package auth

import (
	"testing"
)

func TestHashPassword(t *testing.T) {
	password := "password"

	hash, err := HashPassword(password)

	if err != nil {
		t.Errorf(`HashPassword("password") gave an error: %v`, err)
	}

	same, err := CheckPasswordHash(password, hash)
	if !same || err != nil {
		t.Errorf(`HashPassword("password") = %v, nil\nCheckPasswordHash("password", %v) = %v, %v, want match for true, nil`, hash, hash, same, err)
	}
}
