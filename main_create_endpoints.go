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

	state.mux.HandleFunc("POST /api/users", state.handlerApiCreateUser)
	state.mux.HandleFunc("PUT /api/users", state.handlerApiUpdateUser)
	state.mux.HandleFunc("POST /api/login", state.handlerApiLogin)
	state.mux.HandleFunc("POST /api/refresh", state.handlerApiRefreshToken)
	state.mux.HandleFunc("POST /api/revoke", state.handlerApiRevokeToken)

	state.mux.HandleFunc("POST /api/chirps", state.handlerApiCreateChirp)
	state.mux.HandleFunc("GET /api/chirps", state.handlerApiGetAllChirps)
	state.mux.HandleFunc("GET /api/chirps/{chirpID}", state.handlerGetSingleChirp)
	state.mux.HandleFunc("DELETE /api/chirps/{chirpID}", state.handlerApiDeleteChirp)

	// ----------- Admin Handlers ----------------

	state.mux.HandleFunc("GET /admin/metrics", state.metricsHandler)
	state.mux.HandleFunc("POST /admin/reset", state.handlerApiReset)

	// ----------- Webhooks ----------------

	state.mux.HandleFunc("POST /api/polka/webhooks", state.handlerWebhookPolka)

	return nil
}
