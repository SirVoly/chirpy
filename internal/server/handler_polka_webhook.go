package server

import (
	"encoding/json"
	"github/SirVoly/chirpy/internal/auth"
	"net/http"

	"github.com/google/uuid"
)

// POST /api/polka/webhooks

func (cfg *apiConfig) handlerPolkaWebhooks(w http.ResponseWriter, r *http.Request) {

	// Authentication
	token, err := auth.GetAPIKey(r.Header)
	if (err != nil) || (token != cfg.PolkaKey) {
		respondWithError(w, http.StatusUnauthorized, "", err)
		return
	}

	type parameters struct {
		Event string `json:"event"`
		Data  struct {
			UserID string `json:"user_id"`
		} `json:"data"`
	}
	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err = decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error decoding parameters", err)
		return
	}

	if params.Event != "user.upgraded" {
		respondBasicString(w, http.StatusNoContent, "")
	}

	user_id, err := uuid.Parse(params.Data.UserID)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "", err)
		return
	}

	err = cfg.db.UpgradeUserToChirpyRed(r.Context(), user_id)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "", err)
		return
	}

	respondBasicString(w, http.StatusNoContent, "")

}
