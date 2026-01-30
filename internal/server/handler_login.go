package server

import (
	"encoding/json"
	"fmt"
	"github/SirVoly/chirpy/internal/auth"
	"github/SirVoly/chirpy/internal/database"
	"net/http"
	"strconv"
	"time"
)


	const tokenExpiresInSeconds = 3600
	const refreshTokenExpiresInSeconds = 3600 * 24 * 60 // 60 Days

// "POST /api/login"
func (cfg *apiConfig) handlerLogin(w http.ResponseWriter, r *http.Request) {

	type parameters struct {
		Email            string `json:"email"`
		Password         string `json:"password"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error decoding parameters", err)
		return
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

	token, err := auth.MakeToken(usr.ID, cfg.JWTsecret, tokenExpiresInSeconds)

	refreshToken, err := auth.MakeRefreshToken()
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error creating refresh token", err)
		return
	}
	dur, _ := time.ParseDuration(fmt.Sprintf("%s%s", strconv.Itoa(refreshTokenExpiresInSeconds), "s"))
	expireTime := time.Now().Add(dur)

	_, err = cfg.db.CreateRefreshToken(r.Context(), database.CreateRefreshTokenParams{
		Token: refreshToken,
		UserID: usr.ID,
		ExpiresAt: expireTime,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error creating refresh token", err)
		return
	}

	respondWithJSON(w, http.StatusOK, createJSONLoginUser(usr, token, refreshToken))
}
