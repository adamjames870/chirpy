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

type paramsUpdateUser struct {
	id       uuid.UUID
	Password string `jason:"password"`
	Email    string `json:"email"`
}

func (s *apiState) handlerApiUpdateUser(w http.ResponseWriter, r *http.Request) {

	// PUT api/users
	// header to contain access token
	// user id extracted from token
	// email and password updated from params

	decoder := json.NewDecoder(r.Body)
	paramsUser := paramsUpdateUser{}
	errDecode := decoder.Decode(&paramsUser)
	if errDecode != nil {
		respondWithError(w, 400, "unable to decode json: "+errDecode.Error())
		return
	}

	if paramsUser.Email == "" || paramsUser.Password == "" {
		respondWithError(w, 400, "email and password not provided")
		return
	}

	inputTkn, errTkn := auth.GetBearerToken(r.Header)
	if errTkn != nil {
		respondWithError(w, 401, "unable to parse input token: "+errTkn.Error())
		return
	}

	usrId, errUsrId := auth.ValidateJWT(inputTkn, s.secret_string)
	if errUsrId != nil {
		respondWithError(w, 401, errUsrId.Error())
		return
	}

	paramsUser.id = usrId
	updatedUser, errUpdate := updateUser(r.Context(), *s.dbQueries, paramsUser)
	if errUpdate != nil {
		respondWithError(w, 400, "unable to create user: "+errUpdate.Error())
		return
	}

	rv := user{
		Id:        usrId,
		CreatedAt: updatedUser.CreatedAt,
		UpdatedAt: updatedUser.UpdatedAt,
		Email:     updatedUser.Email,
	}

	respondWithJSON(w, 200, rv)

}

func updateUser(ctx context.Context, db database.Queries, params paramsUpdateUser) (database.User, error) {
	pwd, errPwd := auth.HashPassword(params.Password)
	if errPwd != nil {
		return database.User{}, errors.New("unable to create hash: " + errPwd.Error())
	}

	newParams := database.UpdateUserNameAndEmailParams{
		ID:             params.id,
		Email:          params.Email,
		HashedPassword: pwd,
		UpdatedAt:      time.Now(),
	}

	return db.UpdateUserNameAndEmail(ctx, newParams)

}
