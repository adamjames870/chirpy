package main

import (
	"net/http"
)

func (s *apiState) handlerApiGetAllChirps(w http.ResponseWriter, r *http.Request) {

	// GET /api/chirps

	dbChirps, errChirps := s.dbQueries.GetAllChirps(r.Context())
	if errChirps != nil {
		respondWithError(w, 400, "unable to fetch chirps: "+errChirps.Error())
	}

	countChirps := len(dbChirps)
	rv := make([]chirp, 0, countChirps)
	for _, dbChirp := range dbChirps {
		rv = append(rv, chirp{
			Id:        dbChirp.ID,
			CreatedAt: dbChirp.CreatedAt,
			UpdatedAt: dbChirp.UpdatedAt,
			Body:      dbChirp.Body,
			UserId:    dbChirp.UserID,
		})
	}

	respondWithJSON(w, 200, rv)

}
