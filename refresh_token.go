package main

import (
	"log"
	"net/http"
	"time"

	"github.com/Nischaldh/Chirpy/internal/auth"
)

func (cfg *apiConfig) refresh(w http.ResponseWriter, r *http.Request) {
	refreshToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Unauthorized", nil)
		return
	}

	user, err := cfg.db.GetUserFromRefreshToken(
		r.Context(),
		refreshToken,
	)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Unauthorized", nil)
		return
	}

	token, err := auth.MakeJWT(
		user.ID,
		cfg.jwtSecret,
		time.Hour,
	)
	if err != nil {
		log.Printf("Error creating access token: %v", err)
		respondWithError(w, 500, "Something went wrong", nil)
		return
	}

	respondWithJSON(w, http.StatusOK, map[string]string{
		"token": token,
	})
}


func (cfg *apiConfig) revoke(w http.ResponseWriter, r *http.Request) {
	refreshToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Unauthorized", nil)
		return
	}

	err = cfg.db.RevokeRefreshToken(
		r.Context(),
		refreshToken,
	)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Unauthorized", nil)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}