package main

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/adamjames870/chirpy/internal/database"
	"github.com/google/uuid"
)

func (s *apiState) handlerApiCreateUser(w http.ResponseWriter, r *http.Request) {

	// POST api/users

	type parameters struct {
		Email string `json:"email"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	errDecode := decoder.Decode(&params)
	if errDecode != nil {
		respondWithError(w, 400, "unable to decode json: "+errDecode.Error())
		return
	}

	usr, errUsr := writeNewUser(r.Context(), *s.dbQueries, params.Email)
	if errUsr != nil {
		respondWithError(w, 400, "unable to create user: "+errUsr.Error())
		return
	}

	rv := user{
		Id:        usr.ID,
		CreatedAt: usr.CreatedAt,
		UpdatedAt: usr.UpdatedAt,
		Email:     usr.Email,
	}

	respondWithJSON(w, 201, rv)

}

func writeNewUser(ctx context.Context, db database.Queries, email string) (database.User, error) {
	params := database.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Email:     email,
	}
	return db.CreateUser(ctx, params)
}
