package server

import (
	"database/sql"
	"github/SirVoly/chirpy/internal/database"
	"log"
	"net/http"
	"os"
	"sync/atomic"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func Run() {
	const filepathRoot = "."
	const port = "8080"

	godotenv.Load()
	dbURL := os.Getenv("DB_URL")
	if dbURL == "" {
		log.Fatal("DB_URL must be set")
	}
	platform := os.Getenv("PLATFORM")
	if platform == "" {
		log.Fatal("PLATFORM must be set")
	}
	secret := os.Getenv("JWTSECRET")
	if secret == "" {
		log.Fatal("JWTSECRET must be set")
	}


	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatalf("Error opening database %s", err)
	}

	dbQueries := database.New(db)

	apiCfg := apiConfig{
		fileserverHits: atomic.Int32{},
		db:             dbQueries,
		platform:       platform,
		JWTsecret:			secret,
	}

	serverMux := http.NewServeMux()

	// File Server handler
	serverMux.Handle("/app/", apiCfg.middlewareMetricsInc(middlewareLog(http.StripPrefix("/app/", http.FileServer(http.Dir(filepathRoot))))))

	serverMux.HandleFunc("GET /api/healthz", handlerHealthz)
	serverMux.HandleFunc("GET /admin/metrics", apiCfg.showMetricsHandler)

	serverMux.HandleFunc("POST /api/users", apiCfg.handlerUsersCreate)
	serverMux.HandleFunc("POST /api/login", apiCfg.handlerLogin)

	serverMux.HandleFunc("POST /api/chirps", apiCfg.handlerChirpsCreate)
	serverMux.HandleFunc("GET /api/chirps", apiCfg.handlerChirpsRetrieveAll)
	serverMux.HandleFunc("GET /api/chirps/{chirpID}", apiCfg.handlerChirpsGet)

	serverMux.HandleFunc("POST /admin/reset", apiCfg.handlerReset)



	server := http.Server{
		Handler: serverMux,
		Addr:    ":" + port,
	}

	server.ListenAndServe()
}
