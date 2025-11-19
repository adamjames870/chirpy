package main

import (
	"fmt"
	"net/http"
	"sync/atomic"
)

type apiConfig struct {
	fileserverHits atomic.Int32
}

func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.fileserverHits.Add(1)
		next.ServeHTTP(w, r)
	})
}

func main() {
	var apiCfg apiConfig
	mux := http.NewServeMux()

	// ----------- File Handlers ---------------

	appHandler := http.StripPrefix("/app", http.FileServer(http.Dir(".")))
	mux.Handle("/app/", apiCfg.middlewareMetricsInc(appHandler))

	assetHandler := http.StripPrefix("/assets/", http.FileServer(http.Dir("./assets/")))
	mux.Handle("/assets/", assetHandler)

	// ----------- API Handlers ----------------

	mux.HandleFunc("GET /api/healthz", readinessHandler)

	// ----------- Admin Handlers ----------------

	mux.HandleFunc("GET /admin/metrics", apiCfg.metricsHandler)
	mux.HandleFunc("POST /admin/reset", apiCfg.metricsResetHandler)

	server := http.Server{
		Handler: mux,
		Addr:    ":8080",
	}
	server.ListenAndServe()
}

func readinessHandler(writer http.ResponseWriter, req *http.Request) {
	writer.Header().Set("Content-Type", "text/plain; charset=utf-8")
	writer.WriteHeader(http.StatusOK)
	writer.Write([]byte("OK\n"))
}

func (cfg *apiConfig) metricsHandler(writer http.ResponseWriter, req *http.Request) {
	writer.Header().Set("Content-Type", "text/html; charset=utf-8")
	writer.WriteHeader(http.StatusOK)
	output := fmt.Sprintf("<html>\n<body>\n<h1>Welcome, Chirpy Admin!</h1><p>Chirpy has been visited %d times!</p></body></html>",
		cfg.fileserverHits.Load())
	writer.Write([]byte(output))
}

func (cfg *apiConfig) metricsResetHandler(writer http.ResponseWriter, req *http.Request) {
	writer.Header().Set("Content-Type", "text/plain; charset=utf-8")
	writer.WriteHeader(http.StatusOK)
	cfg.fileserverHits.Store(0)
	writer.Write([]byte("Reset metrics to zero"))
}
