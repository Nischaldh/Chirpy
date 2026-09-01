package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/Nischaldh/Chirpy/internal/auth"
	"github.com/Nischaldh/Chirpy/internal/database"
	"github.com/google/uuid"
)

type User struct {
	ID           uuid.UUID `json:"id"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	Email        string    `json:"email"`
	Token        string    `json:"token,omitempty"`
	RefreshToken string    `json:"refresh_token,omitempty"`
}

func (cfg *apiConfig) createUsers(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	type ReqBody struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	req := ReqBody{}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("Error decoding paramters: %s", err)
		respondWithError(w, 500, "Something went wrong", err)
		return
	}
	if req.Email == "" || req.Password == "" {
		respondWithError(w, http.StatusBadRequest, "Email and Password both are required", nil)
		return
	}
	hash, err := auth.HashPassword(req.Password)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error while hashing the password", nil)
		return
	}
	createdUser, err := cfg.db.CrateUser(r.Context(), database.CrateUserParams{
		Email:    req.Email,
		Password: hash,
	})
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

func (cfg *apiConfig) loginUsers(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	type ReqBody struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	req := ReqBody{}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("Error decoding paramters: %s", err)
		respondWithError(w, 500, "Something went wrong", err)
		return
	}
	if req.Email == "" || req.Password == "" {
		respondWithError(w, http.StatusBadRequest, "Email and Password both are required", nil)
		return
	}
	u, err := cfg.db.GetUserByEmail(r.Context(), req.Email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			respondWithError(w, 404, "Cannot find the user", nil)
			return
		}

		respondWithError(w, 500, "Couldn't get user", err)
		return
	}
	match, err := auth.CheckPasswordHash(req.Password, u.Password)
	if !match {
		respondWithError(w, http.StatusUnauthorized, "Invalid Credentials", nil)
		return
	}
	if err != nil {
		log.Printf("Error decoding paramters: %s", err)
		respondWithError(w, 500, "Something went wrong", err)
		return
	}
	expiresIn := time.Duration(3600) * time.Second
	token, err := auth.MakeJWT(u.ID, cfg.jwtSecret, expiresIn)
	if err != nil {
		log.Printf("Error decoding paramters: %s", err)
		respondWithError(w, 500, "Something went wrong", err)
		return
	}
	refreshToken := auth.MakeRefreshToken()
	_, err = cfg.db.CreateRefreshToken(
		r.Context(),
		database.CreateRefreshTokenParams{
			Token:  refreshToken,
			UserID: u.ID,
		},
	)
	if err != nil {
		log.Printf("Error creating refresh token: %v", err)
		respondWithError(w, 500, "Something went wrong", err)
		return
	}
	respondWithJSON(w, 200, User{
		ID:           u.ID,
		CreatedAt:    u.CreatedAt,
		UpdatedAt:    u.UpdatedAt,
		Email:        u.Email,
		Token:        token,
		RefreshToken: refreshToken,
	})
}
