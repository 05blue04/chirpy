package main

import (
	"context"
	"encoding/json"
	"net/http"
	"sort"
	"strings"
	"time"

	"github.com/05blue04/chirpy/internal/auth"
	"github.com/05blue04/chirpy/internal/database"
	"github.com/google/uuid"
)

type Chirp struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Body      string    `json:"body"`
	UserID    uuid.UUID `json:"user_id"`
}

func (cfg *apiConfig) createChirpHandler(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Body string `json:"body"`
	}

	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "error extracting bearer from request", err)
		return
	}

	userID, err := auth.ValidateJWT(token, cfg.secret)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "unable to grant access", err)
		return
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}

	err = decoder.Decode(&params)
	if err != nil {
		respondWithError(w, 400, "Error decoding body", err)
		return
	}

	if len(params.Body) > 140 {
		respondWithError(w, 400, "Chirp is too long", nil)
		return
	}

	blocked := map[string]struct{}{
		"kerfuffle": {},
		"sharbert":  {},
		"fornax":    {},
	}
	clean := cleanString(params.Body, blocked)

	c, err := cfg.db.CreateChirp(context.Background(), database.CreateChirpParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Body:      clean,
		UserID:    userID,
	})

	if err != nil {
		respondWithError(w, 400, "error creating chirp", err)
		return
	}

	respondWithJSON(w, 201, Chirp{
		ID:        c.ID,
		CreatedAt: c.CreatedAt,
		UpdatedAt: c.UpdatedAt,
		Body:      c.Body,
		UserID:    c.UserID,
	})

}

func cleanString(body string, blocked map[string]struct{}) string {
	check := strings.Fields(body)

	for i, str := range check { // using a map to create a set to access the words is more optimal in this case but whatever
		str := strings.ToLower(str)
		_, ok := blocked[str]
		if ok {
			check[i] = "****"
		}
	}

	return strings.Join(check, " ")
}

func (cfg *apiConfig) getChirpsHandler(w http.ResponseWriter, r *http.Request) {

	var chirps []database.Chirp
	var err error

	authIDStr := r.URL.Query().Get("author_id")
	sortParam := r.URL.Query().Get("sort")
	if authIDStr != "" {
		err := uuid.Validate(authIDStr)
		if err != nil {
			respondWithError(w, http.StatusBadRequest, "invalid author_id uuid", err)
			return
		}

		authID, err := uuid.Parse(authIDStr)
		if err != nil {
			respondWithError(w, http.StatusBadRequest, "error parsing author_id", err)
			return
		}

		chirps, err = cfg.db.GetChirpsByAuthID(r.Context(), authID)
		if err != nil {
			respondWithError(w, 500, "error getting chirps by authID", err)
			return
		}

	} else {
		chirps, err = cfg.db.GetChirps(r.Context())
		if err != nil {
			respondWithError(w, 500, "error getting chirps from db", err)
			return
		}
	}

	jsonChirps := make([]Chirp, len(chirps))
	for i, c := range chirps {
		jsonChirps[i] = Chirp{
			ID:        c.ID,
			CreatedAt: c.CreatedAt,
			UpdatedAt: c.UpdatedAt,
			Body:      c.Body,
			UserID:    c.UserID,
		}
	}

	if sortParam == "desc" {
		sort.Slice(jsonChirps, func(i, j int) bool {
			return jsonChirps[i].CreatedAt.After(jsonChirps[j].CreatedAt)
		})
	}

	respondWithJSON(w, http.StatusOK, jsonChirps)
}

func (cfg *apiConfig) getChirpByIDHandler(w http.ResponseWriter, r *http.Request) {
	chirpID := r.PathValue("chirpID")

	err := uuid.Validate(chirpID)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid uuid in request", err)
		return
	}

	chirpUUID, err := uuid.Parse(chirpID)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "error converting id to uuid", err)
		return
	}

	c, err := cfg.db.GetChirpById(context.Background(), chirpUUID)
	if err != nil {
		respondWithError(w, 404, "Chirp with requested id does not exist", err)
		return
	}

	respondWithJSON(w, http.StatusOK, Chirp{
		ID:        c.ID,
		CreatedAt: c.CreatedAt,
		UpdatedAt: c.UpdatedAt,
		Body:      c.Body,
		UserID:    c.UserID,
	})

}

func (cfg *apiConfig) deleteChirpHandler(w http.ResponseWriter, r *http.Request) {
	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "error extracting bearer from request", err)
		return
	}

	userID, err := auth.ValidateJWT(token, cfg.secret)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "unable to grant access", err)
		return
	}

	chirpIDstring := r.PathValue("chirpID")

	err = uuid.Validate(chirpIDstring)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid uuid in request", err)
		return
	}

	chirpID, err := uuid.Parse(chirpIDstring)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "error converting id to uuid", err)
		return
	}

	chirp, err := cfg.db.GetChirpById(r.Context(), chirpID)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "invalid chirpID", err)
		return
	}

	if chirp.UserID != userID {
		respondWithError(w, http.StatusForbidden, "Only users the author of this chirp can delete", err)
		return
	}

	err = cfg.db.DeleteChirpByID(r.Context(), chirpID)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "invalid chirpID", err)
		return
	}

	w.WriteHeader(http.StatusNoContent)

}
