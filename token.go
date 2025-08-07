package main

import (
	"context"
	"net/http"
	"time"

	"github.com/05blue04/chirpy/internal/auth"
)

func (cfg *apiConfig) refreshHandler(w http.ResponseWriter, r *http.Request) {

	type response struct {
		Token string `json:"token"`
	}

	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, 401, "unable to get token from header", err)
		return
	}

	refreshToken, err := cfg.db.GetTokenByID(context.Background(), token)

	if err != nil {
		respondWithError(w, 401, "not authorized boy", err)
		return
	}

	if time.Now().After(refreshToken.ExpiresAt) {
		respondWithError(w, 401, "refresh token has expired", nil)
		return
	}

	accessToken, err := auth.MakeJWT(refreshToken.UserID, cfg.secret, time.Hour)
	if err != nil {
		respondWithError(w, 500, "error creating New access Token", err)
		return
	}

	respondWithJSON(w, 200, response{
		Token: accessToken,
	})

}

func (cfg *apiConfig) revokeHandler(w http.ResponseWriter, r *http.Request) {
	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, 401, "unable to get token from header", err)
		return
	}

	err = cfg.db.RevokeToken(context.Background(), token)

	if err != nil {
		respondWithError(w, 401, "The token provided doesn't match our records", err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
