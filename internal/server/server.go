package server

import (
	"database/sql"
	"encoding/json"
	"github/SirVoly/chirpy/internal/database"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/google/uuid"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func Run() {
	apiCfg := apiConfig{}

	godotenv.Load()
	dbURL := os.Getenv("DB_URL")
	apiCfg.platform = os.Getenv("PLATFORM")
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		return
	}
	
	apiCfg.db = database.New(db)

	serverMux := http.NewServeMux()

	// A link with the suffix /app/ will be redirected to this handler, which in turn will remove the /app/ prefix, so the FileServer can interpret this in the same way as the files.
	serverMux.Handle("/app/", apiCfg.middlewareMetricsInc(middlewareLog(http.StripPrefix("/app/", http.FileServer(http.Dir("."))))))

	serverMux.HandleFunc("GET /api/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	serverMux.HandleFunc("GET /api/chirps", func(w http.ResponseWriter, r *http.Request) {
		chirps, err := apiCfg.db.ListChirps(r.Context())
		if err != nil {
			log.Printf("Error getting chirps: %s", err)
			w.WriteHeader(500)
			w.Write([]byte("Internal Server Error"))
			return
		}

		listChirps := make([]JSON_Chirp, 0)

		for _, c := range chirps {
			listChirps = append(listChirps, createJSONChirp(c))
		}

		respondWithJSON(w, http.StatusOK, listChirps)
	})

	serverMux.HandleFunc("POST /api/chirps", func(w http.ResponseWriter, r *http.Request) {

		type parameters struct {
			Body   string `json:"body"`
			UserID string `json:"user_id"`
		}

		decoder := json.NewDecoder(r.Body)
		params := parameters{}
		err := decoder.Decode(&params)
		if err != nil {
			log.Printf("Error decoding parameters: %s", err)
			w.WriteHeader(500)
			w.Write([]byte("Internal Server Error"))
			return
		}
		user_id := uuid.MustParse(params.UserID)

		msg := params.Body

		// Validate Chirp
		valid := len(msg) <= 140
		if !valid {
			respondWithError(w, 400, "Chirp is too long")
			return
		}
		msg = cleanChirp(msg)

		// Upload Chirp
		chirp, err := apiCfg.db.CreateChirp(r.Context(), database.CreateChirpParams{
			Body:   msg,
			UserID: user_id,
		})

		if err != nil {
			respondWithError(w, http.StatusInternalServerError, err.Error())
			return
		}

		respondWithJSON(w, http.StatusCreated, createJSONChirp(chirp))
	})

	serverMux.HandleFunc("GET /admin/metrics", apiCfg.showMetricsHandler)
	serverMux.HandleFunc("POST /admin/reset", func(w http.ResponseWriter, r *http.Request) {
		if apiCfg.platform != "dev" {
			w.WriteHeader(http.StatusForbidden)
			w.Write([]byte("Forbidden\n"))
			return
		}

		err := apiCfg.db.DeleteAllUsers(r.Context())
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Failed to delete all users\n"))
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Reset all users\n"))
	})

	serverMux.HandleFunc("GET /api/chirps/{chirpID}", func(w http.ResponseWriter, r *http.Request) {

		chirp_id := r.PathValue("chirpID")

		// Get Chirp
		chirp, err := apiCfg.db.GetChirpFromID(r.Context(), uuid.MustParse(chirp_id))
		
		if err != nil {
			respondWithError(w, 404, "Chirp not found")
			return
		}

		respondWithJSON(w, http.StatusOK, createJSONChirp(chirp))
	})


	serverMux.HandleFunc("POST /api/users", func(w http.ResponseWriter, r *http.Request) {

		type parameters struct {
			Email string `json:"email"`
		}

		decoder := json.NewDecoder(r.Body)
		params := parameters{}
		err := decoder.Decode(&params)
		if err != nil {
			log.Printf("Error decoding parameters: %s", err)
			w.WriteHeader(500)
			w.Write([]byte(""))
			return
		}

		// Create user
		usr, err := apiCfg.db.CreateUser(r.Context(), params.Email)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, err.Error())
			return
		}

		respondWithJSON(w, http.StatusCreated, createJSONUser(usr))
	})

	server := http.Server{
		Handler: serverMux,
		Addr:    ":8080",
	}

	server.ListenAndServe()
}

func respondWithError(w http.ResponseWriter, code int, msg string) {
	dat, err := json.Marshal(struct {
		Error string `json:"error"`
	}{Error: msg})
	if err != nil {
		log.Printf("Error marshalling JSON: %s", err)
		w.WriteHeader(500)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(dat)
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {

	dat, err := json.Marshal(payload)
	if err != nil {
		log.Printf("Error marshalling JSON: %s", err)
		w.WriteHeader(500)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(dat)
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
