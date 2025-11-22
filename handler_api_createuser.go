package main

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/adamjames870/chirpy/internal/auth"
	"github.com/adamjames870/chirpy/internal/database"
	"github.com/google/uuid"
)

type paramsCreateUser struct {
	Password string `jason:"password"`
	Email    string `json:"email"`
}

func (s *apiState) handlerApiCreateUser(w http.ResponseWriter, r *http.Request) {

	// POST api/users

	decoder := json.NewDecoder(r.Body)
	params := paramsCreateUser{}
	errDecode := decoder.Decode(&params)
	if errDecode != nil {
		respondWithError(w, 400, "unable to decode json: "+errDecode.Error())
		return
	}

	usr, errUsr := writeNewUser(r.Context(), *s.dbQueries, params)
	if errUsr != nil {
		respondWithError(w, 400, "unable to create user: "+errUsr.Error())
		return
	}

	rv := user{
		Id:          usr.ID,
		CreatedAt:   usr.CreatedAt,
		UpdatedAt:   usr.UpdatedAt,
		Email:       usr.Email,
		IsChirpyRed: usr.IsChirpyRed,
	}

	respondWithJSON(w, 201, rv)

}

func writeNewUser(ctx context.Context, db database.Queries, params paramsCreateUser) (database.User, error) {
	pwd, errPwd := auth.HashPassword(params.Password)
	if errPwd != nil {
		return database.User{}, errors.New("unable to create hash: " + errPwd.Error())
	}

	newParams := database.CreateUserParams{
		ID:             uuid.New(),
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
		Email:          params.Email,
		HashedPassword: pwd,
	}
	return db.CreateUser(ctx, newParams)
}
