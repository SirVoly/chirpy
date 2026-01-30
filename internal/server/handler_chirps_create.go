package server

import (
	"encoding/json"
	"github/SirVoly/chirpy/internal/auth"
	"github/SirVoly/chirpy/internal/database"
	"net/http"
	"strings"
)

// POST /api/chirps
func (cfg *apiConfig) handlerChirpsCreate(w http.ResponseWriter, r *http.Request) {

	type parameters struct {
		Body   string `json:"body"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error decoding parameters", err)
		return
	}

	// Authentication

	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Unauthorized", err)
		return
	}

	user_id, err := auth.ValidateJWT(token, cfg.JWTsecret)

	if (err != nil) {
		respondWithError(w, http.StatusUnauthorized, "Unauthorized", err)
		return
	}

	msg := params.Body

	// Validate Chirp
	if !(len(msg) <= 140) {
		respondWithError(w, http.StatusBadRequest, "Chirp is too long", err)
		return
	}
	msg = cleanChirp(msg)

	// Upload Chirp
	chirp, err := cfg.db.CreateChirp(r.Context(), database.CreateChirpParams{
		Body:   msg,
		UserID: user_id,
	})

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error creating chirp in database.", err)
		return
	}

	respondWithJSON(w, http.StatusCreated, createJSONChirp(chirp))
}

func cleanChirp(input string) string {
	badWords := [3]string{"kerfuffle", "sharbert", "fornax"}

	words := strings.Split(input, " ")
	for index, w := range words {
		for _, b := range badWords {
			if strings.ToLower(w) == b {
				words[index] = "****"
				break
			}
		}
	}
	return strings.Join(words, " ")
}
