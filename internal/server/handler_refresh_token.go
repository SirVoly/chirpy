package server

import (
	"github/SirVoly/chirpy/internal/auth"
	"net/http"
)

// POST /api/refresh
func (cfg *apiConfig) handlerRefresh(w http.ResponseWriter, r *http.Request) {
	// Authentication
	refreshToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Unauthorized", err)
		return
	}

	user, err := cfg.db.GetUserFromValidRefreshToken(r.Context(), refreshToken)

	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Unauthorized", err)
		return
	}

	token, err := auth.MakeToken(user.ID, cfg.JWTsecret, tokenExpiresInSeconds)

	respondWithJSON(w, http.StatusOK, struct {Token         string `json:"token"`}{
		Token: token,
	})
}

func (cfg * apiConfig) handlerRevoke(w http.ResponseWriter, r *http.Request) {
	// Authentication
	refreshToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Unauthorized", err)
		return
	}

	cfg.db.RevokeRefreshToken(r.Context(), refreshToken)
	respondBasicString(w, http.StatusNoContent, "")
}