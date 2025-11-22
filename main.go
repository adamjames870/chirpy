package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func (cfg *apiState) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.fileserverHits.Add(1)
		next.ServeHTTP(w, r)
	})
}

func main() {

	var state apiState
	godotenv.Load()
	state.LoadState()

	state.mux = http.NewServeMux()
	state.CreateEndpoints()

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	server := http.Server{
		Handler: state.mux,
		Addr:    ":" + port,
	}
	server.ListenAndServe()
}

func readinessHandler(writer http.ResponseWriter, req *http.Request) {
	writer.Header().Set("Content-Type", "text/plain; charset=utf-8")
	writer.WriteHeader(http.StatusOK)
	writer.Write([]byte("OK\n"))
}

func (cfg *apiState) metricsHandler(writer http.ResponseWriter, req *http.Request) {
	writer.Header().Set("Content-Type", "text/html; charset=utf-8")
	writer.WriteHeader(http.StatusOK)
	output := fmt.Sprintf("<html>\n<body>\n<h1>Welcome, Chirpy Admin!</h1><p>Chirpy has been visited %d times!</p></body></html>",
		cfg.fileserverHits.Load())
	writer.Write([]byte(output))
}

func (cfg *apiState) metricsResetHandler(writer http.ResponseWriter, req *http.Request) {
	writer.Header().Set("Content-Type", "text/plain; charset=utf-8")
	writer.WriteHeader(http.StatusOK)
	cfg.fileserverHits.Store(0)
	writer.Write([]byte("Reset metrics to zero"))
}
