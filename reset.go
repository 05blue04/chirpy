package main

import (
	"context"
	"net/http"
)

func (cfg *apiConfig) resetHandler(w http.ResponseWriter, r *http.Request) {
	if cfg.platform != "dev" {
		respondWithError(w, http.StatusForbidden, "Must have role dev to hit this endpoint", nil)
		return
	}

	err := cfg.db.ClearUsers(context.Background())
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error clearing users from database", err)
		return
	}
	cfg.fileserverHits.Store(0)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("users table has been cleared and hits have been reset"))
}
