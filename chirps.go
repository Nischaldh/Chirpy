package main

import (
	"database/sql"
	"encoding/json"
	"errors"
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

func (cfg *apiConfig) createChirps (w http.ResponseWriter, r *http.Request){
	defer r.Body.Close()
	type ReqBody struct{
		BODY string `json:"body"`
		USER_ID uuid.UUID `json:"user_id"`
	}
	req:= ReqBody{}
	if err:=json.NewDecoder(r.Body).Decode(&req);err!=nil{
		log.Printf("Error decoding parameters: %s", err)
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
		log.Printf("Error decoding parameters: %s", err)
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

func (cfg *apiConfig) getChirps (w http.ResponseWriter, r *http.Request){
	defer r.Body.Close()
	var chirps []Chirp
	response, err := cfg.db.GetChirps(r.Context())
	if err!=nil{
		log.Printf("Error getting chirps: %s", err)
		respondWithError(w, 500, "Couldn't retrieve chirps", err)
		return
	}
	for _, res:= range response{
		chirps = append(chirps, Chirp{
			ID:res.ID,
			CreatedAt: res.CreatedAt,
			UpdatedAt: res.UpdatedAt,
			Body: res.Body,
			UserID: res.UserID,
		})
	}
	respondWithJSON(w, 200, chirps)
}

func (cfg *apiConfig) getChirp(w http.ResponseWriter, r *http.Request){
	defer r.Body.Close()
	chirpId, err := uuid.Parse(r.PathValue("chirpID"))
	if err!=nil{
		log.Printf("Error while converting chirpId to UUID: %s", err)
		respondWithError(w, 500, "Couldn't convert ChirpId to UUID", err)
		return
	}
	chirp, err := cfg.db.GetChirp(r.Context(),chirpId)
	  if err != nil {
        if errors.Is(err, sql.ErrNoRows) {
            respondWithError(w, 404, "Cannot find the chirp", nil)
            return
        }

        respondWithError(w, 500, "Couldn't get chirp", err)
        return
    }
	respondWithJSON(w, 200, Chirp{
		ID: chirpId,
		UpdatedAt: chirp.UpdatedAt,
		CreatedAt: chirp.CreatedAt,
		Body: chirp.Body,
		UserID: chirp.UserID,
	})
}