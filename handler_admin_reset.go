package main

import "net/http"

func (s *apiState) handlerApiReset(w http.ResponseWriter, r *http.Request) {

	// POST admin/reset

	if s.platform != "dev" {
		respondWithError(w, 403, "")
	}

	s.fileserverHits.Store(0)

	errReset := s.dbQueries.ResetUsers(r.Context())
	if errReset != nil {
		respondWithError(w, 400, "Unable to reset: "+errReset.Error())
		return
	}
	respondWithJSON(w, 200, nil)
}
