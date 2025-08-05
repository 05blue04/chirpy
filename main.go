package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"sync/atomic"

	"github.com/05blue04/chirpy/internal/database"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type apiConfig struct {
	fileserverHits atomic.Int32
	db             *database.Queries
	platform       string
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("error loading environment variables")
	}

	const port = "8080"
	mux := http.NewServeMux()

	db, err := sql.Open("postgres", os.Getenv("DB_URL"))
	if err != nil {
		log.Fatalf("Couldnt connect to database: %v", err)
	}
	cfg := apiConfig{
		db:       database.New(db),
		platform: os.Getenv("PLATFORM"),
	}

	handler := http.StripPrefix("/app/", http.FileServer(http.Dir(".")))

	mux.Handle("/app/", cfg.middlewareMetricsInc(handler))
	mux.Handle("/app/assets/", cfg.middlewareMetricsInc(http.StripPrefix("/app/assets/", http.FileServer(http.Dir("./assets")))))
	mux.HandleFunc("GET /api/healthz", handlerReadiness)
	mux.HandleFunc("GET /admin/metrics", cfg.metricHandler)
	mux.HandleFunc("POST /admin/reset", cfg.resetHandler)
	mux.HandleFunc("POST /api/users", cfg.usersHandler)
	//chirps
	mux.HandleFunc("POST /api/chirps", cfg.createChirpHandler)
	mux.HandleFunc("GET /api/chirps", cfg.getChirpsHandler)
	mux.HandleFunc("GET /api/chirps/{chirpID}", cfg.getChirpByIDHandler)
	server := &http.Server{
		Handler: mux,
		Addr:    ":" + port,
	}

	err = server.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}
