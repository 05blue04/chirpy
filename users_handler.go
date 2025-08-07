package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/05blue04/chirpy/internal/auth"
	"github.com/05blue04/chirpy/internal/database"
	"github.com/google/uuid"
)

type User struct {
	ID            uuid.UUID `json:"id"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
	Email         string    `json:"email"`
	Is_chirpy_red bool      `json:"is_chirpy_red"`
}

func (cfg *apiConfig) usersHandler(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Password string `json:"password"`
		Email    string `json:"email"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}

	err := decoder.Decode(&params)
	if err != nil {
		log.Printf("Error decoding parameters: %s", err)
		respondWithError(w, 400, "couldn't decode parameters", err)
		return
	}

	hash, err := auth.HashPassword(params.Password)
	if err != nil {
		respondWithError(w, 500, "couldn't create hash for password", err)
		return
	}

	u, err := cfg.db.CreateUser(context.Background(), database.CreateUserParams{
		ID:             uuid.New(),
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
		Email:          params.Email,
		HashedPassword: hash,
	})

	if err != nil {
		respondWithError(w, 400, "error creating user", err)
		return
	}

	newUsr := User{
		ID:            u.ID,
		CreatedAt:     u.CreatedAt,
		UpdatedAt:     u.UpdatedAt,
		Email:         u.Email,
		Is_chirpy_red: u.IsChirpyRed,
	}

	respondWithJSON(w, 201, newUsr)
}

func (cfg *apiConfig) loginHandler(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Password           string `json:"password"`
		Email              string `json:"email"`
		Expires_in_seconds int    `json:"expires_in_seconds"`
	}

	type response struct {
		ID            uuid.UUID `json:"id"`
		CreatedAt     time.Time `json:"created_at"`
		UpdatedAt     time.Time `json:"updated_at"`
		Email         string    `json:"email"`
		Token         string    `json:"token"`
		RefreshToken  string    `json:"refresh_token"`
		Is_chirpy_red bool      `json:"is_chirpy_red"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}

	err := decoder.Decode(&params)
	if err != nil {
		log.Printf("Error decoding parameters: %s", err)
		respondWithError(w, 400, "couldn't decode parameters", err)
		return
	}

	u, err := cfg.db.GetUserByEmail(context.Background(), params.Email)
	if err != nil {
		respondWithError(w, 400, "Problem getting user via email", err)
		return
	}

	err = auth.CheckPasswordHash(params.Password, u.HashedPassword)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "password doesn't match our records", err)
		return
	}

	token, err := auth.MakeJWT(u.ID, cfg.secret, time.Hour)
	if err != nil {
		respondWithError(w, 500, "issues generating JWT token for user", err)
		return
	}

	refreshToken := auth.MakeRefreshToken()

	err = cfg.db.CreateToken(context.Background(), database.CreateTokenParams{
		Token:     refreshToken,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		UserID:    u.ID,
		ExpiresAt: time.Now().Add(24 * time.Hour * 60),
		RevokedAt: sql.NullTime{
			Valid: false,
		},
	})

	if err != nil {
		respondWithError(w, 500, "issue generating Refresh token", err)
		return
	}

	respondWithJSON(w, http.StatusOK, response{
		ID:            u.ID,
		CreatedAt:     u.CreatedAt,
		UpdatedAt:     u.UpdatedAt,
		Email:         u.Email,
		Token:         token,
		RefreshToken:  refreshToken,
		Is_chirpy_red: u.IsChirpyRed,
	})
}

func (cfg *apiConfig) updateUserHandler(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Password string `json:"password"`
		Email    string `json:"email"`
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

	hash, err := auth.HashPassword(params.Password)
	if err != nil {
		respondWithError(w, 400, "error creating password hash", err)
		return
	}

	u, err := cfg.db.UpdateUser(context.Background(), database.UpdateUserParams{
		HashedPassword: hash,
		Email:          params.Email,
		ID:             userID,
	})
	if err != nil {
		respondWithError(w, 400, "Error updating user data", err)
		return
	}

	respondWithJSON(w, 200, User{
		ID:            userID,
		CreatedAt:     u.CreatedAt,
		UpdatedAt:     u.UpdatedAt,
		Email:         u.Email,
		Is_chirpy_red: u.IsChirpyRed,
	})

}

func (cfg *apiConfig) polkaHandler(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Event string `json:"event"`
		Data  struct {
			UserID uuid.UUID `json:"user_id"`
		} `json:"data"`
	}

	key, err := auth.GetAPIKey(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "apiKey doesn't match", err)
		return
	}

	if key != cfg.apiKey {
		respondWithError(w, http.StatusUnauthorized, "apiKey doesn't match", err)
		return
	}
	decoder := json.NewDecoder(r.Body)
	params := parameters{}

	err = decoder.Decode(&params)
	if err != nil {
		log.Printf("Error decoding parameters: %s", err)
		respondWithError(w, 400, "couldn't decode parameters", err)
		return
	}

	if params.Event != "user.upgraded" {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	err = cfg.db.UpdateUserToRed(r.Context(), params.Data.UserID)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "unable to find user", err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
