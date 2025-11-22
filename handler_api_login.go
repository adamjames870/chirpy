package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/adamjames870/chirpy/internal/auth"
	"github.com/adamjames870/chirpy/internal/database"
)

func (s *apiState) handlerApiLogin(w http.ResponseWriter, r *http.Request) {

	// POST/api/login

	type paramsLogin struct {
		Password string `jason:"password"`
		Email    string `json:"email"`
	}

	decoder := json.NewDecoder(r.Body)
	params := paramsLogin{}
	errDecode := decoder.Decode(&params)
	if errDecode != nil {
		respondWithError(w, 400, "unable to decode json: "+errDecode.Error())
		return
	}

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

	refreshTkn, _ := auth.MakeRefreshToken()

	timeNow := time.Now()
	rfParams := database.CreateRefreshTokenParams{
		Token:     refreshTkn,
		CreatedAt: timeNow,
		UpdatedAt: timeNow,
		UserID:    usr.ID,
		ExpiresAt: timeNow.Add(expiryTimeRefreshToken),
	}

	savedRefreshToken, errRefreshToken := s.dbQueries.CreateRefreshToken(r.Context(), rfParams)
	if errRefreshToken != nil {
		respondWithError(w, 400, "error saving refresh token"+errRefreshToken.Error())
	}

	accessTkn, errTkn := auth.MakeJWT(usr.ID, s.secret_string, time.Duration(expiryTimeAccesToken))
	if errTkn != nil {
		respondWithError(w, 400, "unable to create token: "+errTkn.Error())
	}

	rv := user{
		Id:           usr.ID,
		CreatedAt:    usr.CreatedAt,
		UpdatedAt:    usr.UpdatedAt,
		Email:        usr.Email,
		Token:        accessTkn,
		RefreshToken: savedRefreshToken.Token,
		IsChirpyRed:  usr.IsChirpyRed,
	}

	respondWithJSON(w, 200, rv)

}
