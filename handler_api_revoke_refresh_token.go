package main

import (
	"net/http"
	"time"

	"github.com/adamjames870/chirpy/internal/auth"
	"github.com/adamjames870/chirpy/internal/database"
)

func (s *apiState) handlerApiRevokeToken(w http.ResponseWriter, r *http.Request) {

	// POST api/revoke

	inputTkn, errTkn := auth.GetBearerToken(r.Header)
	if errTkn != nil {
		respondWithError(w, 400, "unable to parse input token: "+errTkn.Error())
	}

	revokeParams := database.RevokeTokenByTokenParams{
		UpdatedAt: time.Now(),
		Token:     inputTkn,
	}

	errRevoke := s.dbQueries.RevokeTokenByToken(r.Context(), revokeParams)
	if errRevoke != nil {
		respondWithError(w, 400, "unable to revoke token: "+errRevoke.Error())
	}

	respondWithJSON(w, 204, nil)

}
