package server

import (
	"net/http"

	"github.com/google/uuid"
)

// GET /api/chirps/{chirpID}
func (cfg *apiConfig) handlerChirpsGet(w http.ResponseWriter, r *http.Request) {
	chirp_id := r.PathValue("chirpID")

	// Get Chirp
	chirp, err := cfg.db.GetChirpFromID(r.Context(), uuid.MustParse(chirp_id))
	
	if err != nil {
		respondWithError(w, http.StatusNotFound, "Chirp not found", err)
		return
	}

	respondWithJSON(w, http.StatusOK, createJSONChirp(chirp))
}

// GET /api/chirps
func (cfg *apiConfig) handlerChirpsRetrieveAll(w http.ResponseWriter, r *http.Request) {
	chirps, err := cfg.db.ListChirps(r.Context())
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error gettings chirps from db.", err)
		return
	}

	listChirps := make([]JSON_Chirp, 0)

	for _, c := range chirps {
		listChirps = append(listChirps, createJSONChirp(c))
	}

	respondWithJSON(w, http.StatusOK, listChirps)
}