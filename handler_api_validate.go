package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

var badWords = map[string]struct{}{"kerfuffle": {}, "sharbert": {}, "fornax": {}}

const badWordReplacement string = "****"

func handlerApiValidateChirp(w http.ResponseWriter, r *http.Request) {

	type parameters struct {
		Body string `json:"body"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)

	if err != nil {
		respondWithError(w, 400, "unable to decode json: "+err.Error())
		return
	}

	chirp := cleanText(params.Body)
	if len(chirp) > 140 {
		respondWithError(w, 400, fmt.Sprintf("chirp is too long (%v chars)", len(chirp)))
		return
	}

	type returnVal struct {
		Valid       bool   `json:"valid"`
		CleanedBody string `json:"cleaned_body"`
	}
	val := returnVal{
		Valid:       true,
		CleanedBody: chirp,
	}

	respondWithJSON(w, 200, val)

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
