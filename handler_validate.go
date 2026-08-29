package main

import (
	"encoding/json"
	"log"
	"net/http"
)

func validateInput(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	type ReqBody struct {
		Body string `json:"body"`
	}
	type ResBody struct {
		VALID bool `json:"valid"`
	}

	decoder := json.NewDecoder(r.Body)
	req := ReqBody{}
	err := decoder.Decode(&req)
	if err != nil {
		log.Printf("Error decoding parameters: %s", err)
		respondWithError(w, 500, "Something went wrong", err)
		return
	}
	if len(req.Body) > 140 {
		respondWithError(w, 400, "Chirp is too long", err)
		return
	}
	res := ResBody{
		VALID: true,
	}
	respondWithJSON(w, 200, res)
}