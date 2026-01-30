package auth

import (
	"net/http"
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

func TestGetBearerToken(t *testing.T) {
	tests := []struct {
		name	string
		headers	http.Header
		expResult	string
		wantErr	bool
	}{
		{
			name: "Valid Header",
			headers:	http.Header{
				"Authorization": {"Bearer banaan"},
			},
			expResult:	"banaan",
			wantErr:	false,
		},
		{
			name: "Missing Header",
			headers:	http.Header{
			},
			expResult:	"",
			wantErr:	true,
		},
		{
			name: "Wrong Header",
			headers:	http.Header{
				"Authorization": {"Bear banaan"},
			},
			expResult:	"",
			wantErr:	true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			token, err := GetBearerToken(test.headers)

			if (err != nil) != test.wantErr {
				t.Errorf("GetBearerToken() error = %v, wantErr %v", err, test.wantErr)
			}

			if (test.expResult != token) {
				t.Errorf("ValidateJWT() expects %s, got %s", test.expResult, token)
			}
		})
	}
}
