package main

import "net/http"

func (state *apiState) CreateEndpoints() error {

	// ----------- File Handlers ---------------

	appHandler := http.StripPrefix("/app", http.FileServer(http.Dir(".")))
	state.mux.Handle("/app/", state.middlewareMetricsInc(appHandler))

	assetHandler := http.StripPrefix("/assets/", http.FileServer(http.Dir("./assets/")))
	state.mux.Handle("/assets/", assetHandler)

	// ----------- API Handlers ----------------

	state.mux.HandleFunc("GET /api/healthz", readinessHandler)
	state.mux.HandleFunc("POST /api/validate_chirp", handlerApiValidateChirp)
	state.mux.HandleFunc("POST /api/users", state.handlerApiCreateUser)

	// ----------- Admin Handlers ----------------

	state.mux.HandleFunc("GET /admin/metrics", state.metricsHandler)
	state.mux.HandleFunc("POST /admin/reset", state.handlerApiReset)

	return nil
}
