package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/Nischaldh/Chirpy/internal/database"
	"github.com/google/uuid"
)

type Chirp struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Body      string    `json:"body"`
	UserID    uuid.UUID `json:"user_id"`
}


func (cfg *apiConfig) createChirps (w http.ResponseWriter, r *http.Request){
	defer r.Body.Close()
	type ReqBody struct{
		BODY string `json:"body"`
		USER_ID uuid.UUID `json:"user_id"`
	}
	req:= ReqBody{}
	if err:=json.NewDecoder(r.Body).Decode(&req);err!=nil{
		log.Printf("Error decoding parameters: %s", nil)
		respondWithError(w, 500, "Something went wrong", err)
		return
	}
	if (req.BODY == "" || req.USER_ID == uuid.Nil){
		respondWithError(w, 400, "Please provide chirps and user_id", nil)
		return
	}
	if len(req.BODY) > 140 {
		respondWithError(w, 400, "Chirp is too long", nil)
		return
	}

	cleanedBody := filterProfanity(req.BODY)
	
	chirp, err := cfg.db.CreateChirps(r.Context(), database.CreateChirpsParams{
		Body: cleanedBody,
		UserID: req.USER_ID,
	})
	if err!= nil{
		log.Printf("Error decoding parameters: %s", nil)
		respondWithError(w, 500, "Something went wrong", err)
		return
	}
	respondWithJSON(w, 201, Chirp{
		ID: chirp.ID,
		CreatedAt: chirp.CreatedAt,
		UpdatedAt: chirp.UpdatedAt,
		UserID: chirp.UserID,
		Body: chirp.Body,
	})
}



func filterProfanity(body string) (string) {
	profanityWords := map[string]bool{
		"kerfuffle": true,
		"sharbert":  true,
		"fornax":    true,
	}
	words := strings.Fields(body)
	for i, word := range words {
		lower := strings.ToLower(word)
		if profanityWords[lower] {
			words[i] = "****"
		}
	}
	return strings.Join(words, " ")


}
