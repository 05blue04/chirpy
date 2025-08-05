package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
)

func cleanString(body string) string {
	check := strings.Fields(body)

	for i, str := range check { // using a map to create a set to access the words is more optimal in this case but whatever
		str := strings.ToLower(str)
		switch str {
		case "kerfuffle":
			check[i] = "****"
		case "sharbert":
			check[i] = "****"
		case "fornax":
			check[i] = "****"
		}
	}

	return strings.Join(check, " ")

}
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

	clean := cleanString(c.Body)

	respondWithJSON(w, http.StatusOK, struct {
		Cleaned_body string `json:"cleaned_body"`
	}{
		Cleaned_body: clean,
	})

}
