package main

import (
	"net/http"
	"time"

	"github.com/adamjames870/chirpy/internal/auth"
)

type respondToken struct {
	Token string `json:"token"`
}

func (s *apiState) handlerApiRefreshToken(w http.ResponseWriter, r *http.Request) {

	// POST api/refresh

	inputTkn, errTkn := auth.GetBearerToken(r.Header)
	if errTkn != nil {
		respondWithError(w, 400, "unable to parse input token: "+errTkn.Error())
	}

	usr, errUserId := s.dbQueries.GetUserFromRefreshToken(r.Context(), inputTkn)
	if errUserId != nil {
		respondWithError(w, 401, "no such valid token: "+errUserId.Error())
	}

	if usr.RevokedAt.Valid {
		// valid means not null therefore revoked
		respondWithError(w, 401, "token expired")
	}

	accessTkn, errTkn := auth.MakeJWT(usr.UserID, s.secret_string, time.Duration(expiryTimeAccesToken))
	if errTkn != nil {
		respondWithError(w, 400, "unable to create token: "+errTkn.Error())
	}

	rv := respondToken{
		Token: accessTkn,
	}

	respondWithJSON(w, 200, rv)

}
