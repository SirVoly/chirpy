package server

import (
	"encoding/json"
	"github/SirVoly/chirpy/internal/auth"
	"github/SirVoly/chirpy/internal/database"
	"net/http"
)

// "POST /api/users"
func (cfg *apiConfig) handlerUsersCreate(w http.ResponseWriter, r *http.Request) {

	type parameters struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error decoding parameters", err)
		return
	}

	hash, err := auth.HashPassword(params.Password)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error hashing password", err)
		return
	}

	// Create user
	usr, err := cfg.db.CreateUser(r.Context(), database.CreateUserParams{
		Email:          params.Email,
		HashedPassword: hash,
	})

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error creating user in database", err)
		return
	}

	respondWithJSON(w, http.StatusCreated, createJSONUser(usr))
}

// PUT api/users
func (cfg *apiConfig) handlerUsersUpdate(w http.ResponseWriter, r *http.Request) {

	// Authentication
	user_id, authorized := cfg.authenticate(r.Header)
	if !authorized {
		respondWithError(w, http.StatusUnauthorized, "Unauthorized", nil)
	}

	type parameters struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error decoding parameters", err)
		return
	}

	hash, err := auth.HashPassword(params.Password)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error hashing password", err)
		return
	}

	err = cfg.db.UpdateUser(r.Context(), database.UpdateUserParams{
		ID: user_id,
		Email: params.Email,
		HashedPassword: hash,
	})

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error updating user in database", err)
		return
	}

	user, err := cfg.db.GetUser(r.Context(), params.Email)

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error updating user in database", err)
		return
	}

	respondWithJSON(w, http.StatusOK, createJSONUser(user))
}