package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/adamjames870/chirpy/internal/database"
	"github.com/google/uuid"
)

var badWords = map[string]struct{}{
	"kerfuffle": {}, "sharbert": {}, "fornax": {},
}

const badWordReplacement string = "****"

func (s *apiState) handlerApiCreateChirp(w http.ResponseWriter, r *http.Request) {

	// POST api/chirps

	type parameters struct {
		Body   string    `json:"body"`
		UserId uuid.UUID `json:"user_id"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)

	if err != nil {
		respondWithError(w, 400, "unable to decode json: "+err.Error())
		return
	}

	chirpBody := cleanText(params.Body)
	if len(chirpBody) > 140 {
		respondWithError(w, 400, fmt.Sprintf("chirp is too long (%v chars)", len(chirpBody)))
		return
	}

	newChirp := database.CreateChirpParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Body:      chirpBody,
		UserID:    params.UserId,
	}

	savedChirp, errChirp := s.dbQueries.CreateChirp(r.Context(), newChirp)
	if errChirp != nil {
		respondWithError(w, 400, "unable to create chirp: "+errChirp.Error())
		return
	}

	rv := chirp{
		Id:        savedChirp.ID,
		CreatedAt: savedChirp.CreatedAt,
		UpdatedAt: savedChirp.UpdatedAt,
		Body:      savedChirp.Body,
		UserId:    savedChirp.UserID,
	}

	respondWithJSON(w, 201, rv)

}

func cleanText(txt string) string {
	words := strings.Split(txt, " ")
	for i, word := range words {
		if _, ok := badWords[strings.ToLower(word)]; ok {
			words[i] = badWordReplacement
		}
	}
	return strings.Join(words, " ")
}
