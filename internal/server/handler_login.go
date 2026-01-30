package server

import (
	"encoding/json"
	"fmt"
	"github/SirVoly/chirpy/internal/auth"
	"net/http"
	"strconv"
	"time"
)

// "POST /api/login"
func (cfg *apiConfig) handlerLogin(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Email            string `json:"email"`
		Password         string `json:"password"`
		ExpiresInSeconds int    `json:"expires_in_seconds"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error decoding parameters", err)
		return
	}

	if params.ExpiresInSeconds <= 0 || params.ExpiresInSeconds > 3600 {
		params.ExpiresInSeconds = 3600
	}

	//Get User
	usr, err := cfg.db.GetUser(r.Context(), params.Email)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Incorrect email or password", err)
		return
	}

	correctPassword, err := auth.CheckPasswordHash(params.Password, usr.HashedPassword)
	if err != nil || !correctPassword {
		respondWithError(w, http.StatusUnauthorized, "Incorrect email or password", err)
		return
	}

	dur, _ := time.ParseDuration(fmt.Sprintf("%s%s", strconv.Itoa(params.ExpiresInSeconds), "s"))

	token, err := auth.MakeJWT(usr.ID, cfg.JWTsecret, dur)

	respondWithJSON(w, http.StatusOK, createJSONUser(usr, token))
}
