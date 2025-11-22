package main

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
)

type upgradeParams struct {
	Event string        `json:"event"`
	Data  userDataParam `json:"data"`
}

type userDataParam struct {
	UserID string `json:"user_id"`
}

const upgrade_event string = "user.upgraded"

func (s *apiState) handlerWebhookPolka(w http.ResponseWriter, r *http.Request) {

	decoder := json.NewDecoder(r.Body)
	paramsHook := upgradeParams{}
	errDecode := decoder.Decode(&paramsHook)
	if errDecode != nil {
		respondWithError(w, 400, "unable to decode json: "+errDecode.Error())
		return
	}

	if paramsHook.Event != upgrade_event {
		respondWithJSON(w, 204, nil)
		return
	}

	usrId, _ := uuid.Parse(paramsHook.Data.UserID)
	_, errUpdate := s.dbQueries.MakeUserChirpyRed(r.Context(), usrId)
	if errUpdate != nil {
		respondWithError(w, 404, "failed to update red:"+errUpdate.Error())
		return
	}

	respondWithJSON(w, 204, nil)

}
