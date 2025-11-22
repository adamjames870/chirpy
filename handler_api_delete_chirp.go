package main

import (
	"net/http"

	"github.com/adamjames870/chirpy/internal/auth"
	"github.com/google/uuid"
)

func (s *apiState) handlerApiDeleteChirp(w http.ResponseWriter, r *http.Request) {

	// DELETE /api/chirps/{chirpID}

	chirpId, errId := uuid.Parse(r.PathValue("chirpID"))
	if errId != nil {
		respondWithError(w, 400, "cannot parse user chirp id to uuid: "+errId.Error())
		return
	}

	tkn, errTkn := auth.GetBearerToken(r.Header)
	if errTkn != nil {
		respondWithError(w, 401, "unable to parse token: "+errTkn.Error())
		return
	}

	usrId, errUsrId := auth.ValidateJWT(tkn, s.secret_string)
	if errUsrId != nil {
		respondWithError(w, 403, errUsrId.Error())
		return
	}

	dbChirp, errChirp := s.dbQueries.GetSingleChirp(r.Context(), chirpId)
	if errChirp != nil {
		respondWithError(w, 404, "cannot load chirp: "+errChirp.Error())
		return
	}

	if dbChirp.UserID != usrId {
		respondWithError(w, 403, "not owner of this chirp")
		return
	}

	errDelete := s.dbQueries.DeleteSingleChirp(r.Context(), dbChirp.ID)
	if errDelete != nil {
		respondWithError(w, 404, "unable to delete chirp: "+errDelete.Error())
	}

	respondWithJSON(w, 204, nil)

}
