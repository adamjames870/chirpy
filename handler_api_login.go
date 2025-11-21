package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/adamjames870/chirpy/internal/auth"
)

func (s *apiState) handlerApiLogin(w http.ResponseWriter, r *http.Request) {

	// POST/api/login

	type paramsLogin struct {
		Password         string `jason:"password"`
		Email            string `json:"email"`
		ExpiresInSeconds *int   `json:"expires_in_seconds"`
	}

	decoder := json.NewDecoder(r.Body)
	params := paramsLogin{}
	errDecode := decoder.Decode(&params)
	if errDecode != nil {
		respondWithError(w, 400, "unable to decode json: "+errDecode.Error())
		return
	}

	expiryTime := 60 * 60
	if params.ExpiresInSeconds != nil && *params.ExpiresInSeconds < expiryTime {
		expiryTime = *params.ExpiresInSeconds
	}
	expiryTime = expiryTime * int(time.Millisecond)

	usr, errGetUser := s.dbQueries.GetUserByEmail(r.Context(), params.Email)
	if errGetUser != nil {
		respondWithError(w, 401, "Incorrect email or password")
		return
	}

	authed, errCheckPwd := auth.CheckPasswordHash(params.Password, usr.HashedPassword)
	if errCheckPwd != nil {
		respondWithError(w, 400, "unable to check password: "+errCheckPwd.Error())
		return
	}

	if !authed {
		respondWithError(w, 401, "Incorrect email or password")
	}

	tkn, errTkn := auth.MakeJWT(usr.ID, s.secret_string, time.Duration(expiryTime))
	if errTkn != nil {
		respondWithError(w, 400, "unable to create token: "+errTkn.Error())
	}

	rv := user{
		Id:        usr.ID,
		CreatedAt: usr.CreatedAt,
		UpdatedAt: usr.UpdatedAt,
		Email:     usr.Email,
		Token:     tkn,
	}

	respondWithJSON(w, 200, rv)

}
