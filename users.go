package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Email     string    `json:"email"`
}


func (cfg *apiConfig) createUsers(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	type ReqBody struct {
		Email string `json:"email"`
	}
	req := ReqBody{}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("Error decoding paramters: %s", err)
		respondWithError(w, 500, "Something went wrong", err)
		return
	}
	if req.Email == "" {
		respondWithError(w, http.StatusBadRequest, "Email is required", nil)
		return
	}
	createdUser, err := cfg.db.CrateUser(r.Context(), req.Email)
	if err != nil {
		log.Printf("Error decoding paramters: %s", err)
		respondWithError(w, 500, "Something went wrong", err)
		return
	}
	user := User{
		ID:        createdUser.ID,
		CreatedAt: createdUser.CreatedAt,
		UpdatedAt: createdUser.UpdatedAt,
		Email:     createdUser.Email,
	}
	respondWithJSON(w, 201, user)

}
