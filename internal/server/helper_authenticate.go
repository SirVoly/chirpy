package server

import (
	"github/SirVoly/chirpy/internal/auth"
	"net/http"

	"github.com/google/uuid"
)



func (cfg *apiConfig) authenticate(h http.Header) (uuid.UUID, bool) {
	token, err := auth.GetBearerToken(h)
	if err != nil {
		return uuid.UUID{}, false
	}

	user_id, err := auth.ValidateJWT(token, cfg.JWTsecret)

	if (err != nil) {
		return uuid.UUID{}, false
	}

	return user_id, true
}