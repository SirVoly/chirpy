package server

import (
	"net/http"
)

func Run() {
	serverMux := http.NewServeMux()
	apiCfg := apiConfig{}

	// A link with the suffix /app/ will be redirected to this handler, which in turn will remove the /app/ prefix, so the FileServer can interpret this in the same way as the files.
	serverMux.Handle("/app/", apiCfg.middlewareMetricsInc(middlewareLog(http.StripPrefix("/app/", http.FileServer(http.Dir("."))))))

	serverMux.HandleFunc("GET /healthz", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	serverMux.HandleFunc("GET /metrics", apiCfg.showMetricsHandler)
	serverMux.HandleFunc("POST /reset", apiCfg.resetMetricsHandler)

	server := http.Server{
		Handler: serverMux,
		Addr: ":8080",
	}

	server.ListenAndServe()
}