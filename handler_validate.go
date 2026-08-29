package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
)

type ReqBody struct {
	Body string `json:"body"`
}
type ResBody struct {
	Body string `json:"cleaned_body"`
}

func validateInput(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

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
	res, err := filterProfanity(&req)
	if err != nil {
		respondWithError(w, 500, "Something went wrong", err)
		return
	}
	respondWithJSON(w, 200, res)
}

func filterProfanity(body *ReqBody) (ResBody, error) {
	profanityWords := map[string]bool{
		"kerfuffle": true,
		"sharbert":  true,
		"fornax":    true,
	}
	words := strings.Fields(body.Body)
	for i, word := range words {
		lower := strings.ToLower(word)
		if profanityWords[lower] {
			words[i] = "****"
		}
	}
	return ResBody{
		Body: strings.Join(words, " "),
	}, nil

}
