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
		w.WriteHeader(500)
		return
	}

	if len(c.Body) > 140 {
		dat, err := json.Marshal(struct {
			Error string `json:"error"`
		}{
			Error: "Chirp is too long",
		})
		if err != nil {
			log.Printf("Error marshalling JSON: %s", err)
			w.WriteHeader(500)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(400)
		w.Write(dat)
		return
	}

	dat, err := json.Marshal(struct {
		Valid bool `json:"valid"`
	}{
		Valid: true,
	})

	if err != nil {
		log.Printf("Error marshalling JSON: %s", err)
		w.WriteHeader(500)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	w.Write(dat)
}
