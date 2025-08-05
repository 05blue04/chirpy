package main

import (
	"encoding/json"
	"log"
	"net/http"
)

func (cfg *apiConfig) validateHandler(w http.ResponseWriter, r *http.Request) {
	type chirp struct {
		Body string `json:"body"`
	}

	decoder := json.NewDecoder(r.Body)
	c := chirp{}

	err := decoder.Decode(&c)
	if err != nil {
		log.Printf("Error decoding parameters: %s", err)
		respondWithError(w, 400, "couldn't decode parameters", err)
		return
	}

	if len(c.Body) > 140 {
		respondWithError(w, 400, "Chirp is too long", nil)
		return
	}

	respondWithJSON(w, http.StatusOK, struct {
		Valid bool `json:"valid"`
	}{
		Valid: true,
	})

}
