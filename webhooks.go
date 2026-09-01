package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/Nischaldh/Chirpy/internal/auth"
	"github.com/google/uuid"
)

func (cfg *apiConfig) updateUserStatus(w http.ResponseWriter, r *http.Request){
	defer r.Body.Close()
	
	type ReqBody struct{
		Event string `json:"event"`
		Data  struct{
			UserId uuid.UUID `json:"user_id"`

		} `json:"data"`
	}
	apikey, err := auth.GetAPIKey(r.Header)
	if err!=nil{
		respondWithError(w, 401, "Error getting api key", err)
		return
	}
	if apikey != cfg.polkaKey{
		respondWithError(w,401, "Invalid key", nil)
		return
	}

	req := ReqBody{}
	if err:= json.NewDecoder(r.Body).Decode(&req);err!=nil{
		respondWithError(w, 500, "Error while decoding the request", err)
		return
	}
	if req.Event != "user.upgraded"{
		respondWithError(w, 204, "not supported", nil)
		return
	}
	_,err = cfg.db.UpdateUserChirpy(r.Context(), req.Data.UserId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			respondWithError(w, 404, "Cannot find the user", nil)
			return
		}

		respondWithError(w, 500, "Couldn't get user", err)
		return
	}
	w.WriteHeader(http.StatusNoContent)

}