package server

import (
	"net/http"

	"github.com/google/uuid"
)

// DELETE /api/chirps/{chirpID}
func (cfg *apiConfig) handlerChirpDelete(w http.ResponseWriter, r *http.Request) {
	// Authentication
	user_id, authorized := cfg.authenticate(r.Header)
	if !authorized {
		respondWithError(w, http.StatusUnauthorized, "Unauthorized", nil)
		return
	}
	
	chirp_id := r.PathValue("chirpID")

	// Get Chirp
	chirp, err := cfg.db.GetChirpFromID(r.Context(), uuid.MustParse(chirp_id))
	
	if err != nil {
		respondWithError(w, http.StatusNotFound, "Chirp not found", err)
		return
	}

	if chirp.UserID != user_id {
		respondWithError(w, http.StatusForbidden, "Unauthorized", nil)
		return
	}

	cfg.db.DeleteChirp(r.Context(), chirp.ID)

	respondBasicString(w, http.StatusNoContent, "")
}