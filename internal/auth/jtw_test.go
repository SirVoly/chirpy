package auth

import (
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestValidateJWT(t *testing.T) {
	tests := []struct {
		name        string
		userID      uuid.UUID
		expiresIn   string
		firstSecret string
		lastSecret  string
		wantErr     bool
	}{
		{
			name:        "Valid Token",
			userID:      uuid.New(),
			expiresIn:   "1h",
			firstSecret: "banaan",
			lastSecret:  "banaan",
			wantErr:     false,
		},
		{
			name:        "Expired Token",
			userID:      uuid.New(),
			expiresIn:   "1ms",
			firstSecret: "banaan",
			lastSecret:  "banaan",
			wantErr:     true,
		},
		{
			name:        "Invalid Secret",
			userID:      uuid.New(),
			expiresIn:   "1ms",
			firstSecret: "banaan",
			lastSecret:  "wrongBanaan",
			wantErr:     true,
		},
	}

	for _, c := range tests {
		t.Run(c.name, func(t *testing.T) {
			dur, _ := time.ParseDuration(c.expiresIn)
			tokenString, err := MakeJWT(c.userID, c.firstSecret, dur)

			if err != nil {
				t.Errorf("Unexpected MakeJWT error %v", err)
			}

			sleep, _ := time.ParseDuration("1s")
			time.Sleep(sleep)

			userID, err := ValidateJWT(tokenString, c.lastSecret)

			if (err != nil) != c.wantErr {
				t.Errorf("ValidateJWT() error = %v, wantErr %v", err, c.wantErr)
			}

			if (userID != c.userID) != c.wantErr {
				t.Errorf("ValidateJWT() expects %s, got %s", c.userID.String(), userID.String())
			}
		})

	}
}
