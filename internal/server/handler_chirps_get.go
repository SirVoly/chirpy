package server

import (
	"github/SirVoly/chirpy/internal/database"
	"net/http"
	"sort"

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
// GET /api/chirps?author_id={user_id}
func (cfg *apiConfig) handlerChirpsRetrieveAll(w http.ResponseWriter, r *http.Request) {
	user_id := r.URL.Query().Get("author_id")

	var chirps []database.Chirp
	var err error
	if user_id == "" {
		chirps, err = cfg.db.ListChirps(r.Context())
	} else {
		user_uuid, err := uuid.Parse(user_id)
		if err != nil {
			respondWithError(w, http.StatusBadRequest, "Invalid user ID", err)
			return
		}
		chirps, err = cfg.db.ListChirpsFromUser(r.Context(), user_uuid)
	}

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error gettings chirps from db.", err)
		return
	}

	listChirps := make([]JSON_Chirp, 0)

	// Sort defaults to ASC
	if r.URL.Query().Get("sort") == "desc" {
		sort.Slice(chirps, func(i, j int) bool {
			return chirps[i].CreatedAt.After(chirps[j].CreatedAt)
		})
	}

	for _, c := range chirps {
		listChirps = append(listChirps, createJSONChirp(c))
	}

	respondWithJSON(w, http.StatusOK, listChirps)
}
