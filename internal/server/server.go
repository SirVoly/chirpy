package server

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
)

func Run() {
	serverMux := http.NewServeMux()
	apiCfg := apiConfig{}

	// A link with the suffix /app/ will be redirected to this handler, which in turn will remove the /app/ prefix, so the FileServer can interpret this in the same way as the files.
	serverMux.Handle("/app/", apiCfg.middlewareMetricsInc(middlewareLog(http.StripPrefix("/app/", http.FileServer(http.Dir("."))))))

	serverMux.HandleFunc("GET /api/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	serverMux.HandleFunc("POST /api/validate_chirp", func(w http.ResponseWriter, r *http.Request) {
		
		type parameters struct {
			Body string `json:"body"`
		}

		decoder := json.NewDecoder(r.Body)
		params := parameters{}
		err := decoder.Decode(&params)
		if err != nil {
			log.Printf("Error decoding parameters: %s", err)
			w.WriteHeader(500)
			return
		}

		chirp := params.Body

		// Check if the message is <=140
		valid := len(chirp) <= 140
		if !valid {
			respondWithError(w, 400, "Chirp is too long")
			return
		}
		chirp = cleanChirp(chirp)

		respondWithJSON(w, 200, struct{CleanedBody string `json:"cleaned_body"`}{CleanedBody:chirp})
	})

	serverMux.HandleFunc("GET /admin/metrics", apiCfg.showMetricsHandler)
	serverMux.HandleFunc("POST /admin/reset", apiCfg.resetMetricsHandler)

	server := http.Server{
		Handler: serverMux,
		Addr: ":8080",
	}

	server.ListenAndServe()
}

func respondWithError(w http.ResponseWriter, code int, msg string) {
		dat, err := json.Marshal(struct{Error string `json:"error"`}{Error: msg})
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
	for index, w := range(words) {
		for _, b := range(badWords) {
			if strings.ToLower(w) == b {
				words[index] = "****"
				break
			}
		}
	}
	return strings.Join(words, " ")
}