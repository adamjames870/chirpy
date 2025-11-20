package main

import (
	"net/http"

	"github.com/google/uuid"
)

func (s *apiState) handlerGetSingleChirp(w http.ResponseWriter, r *http.Request) {

	// GET api/chirps/{chirpID}

	chirpId, errId := uuid.Parse(r.PathValue("chirpID"))
	if errId != nil {
		respondWithError(w, 400, "cannot parse user chirp id to uuid: "+errId.Error())
	}

	dbChirp, errChirp := s.dbQueries.GetSingleChirp(r.Context(), chirpId)
	if errChirp != nil {
		respondWithError(w, 404, "cannot load chirp: "+errChirp.Error())
	}

	rv := chirp{
		Id:        dbChirp.ID,
		CreatedAt: dbChirp.CreatedAt,
		UpdatedAt: dbChirp.UpdatedAt,
		Body:      dbChirp.Body,
		UserId:    dbChirp.UserID,
	}

	respondWithJSON(w, 200, rv)

}
