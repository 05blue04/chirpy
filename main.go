package main

import (
	"fmt"
	"log"
	"net/http"
	"sync/atomic"
)

type apiConfig struct {
	fileserverHits atomic.Int32
}

func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.fileserverHits.Add(1)
		next.ServeHTTP(w, r)
	})
}

func (cfg *apiConfig) metricHandler(w http.ResponseWriter, r *http.Request) {
	s := fmt.Sprintf("Hits: %d", cfg.fileserverHits.Load())
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(s))
}
func (cfg *apiConfig) resetHandler(w http.ResponseWriter, r *http.Request) {
	cfg.fileserverHits.Store(0)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("file server hits have been reset to 0"))
}

func main() {
	const port = "8080"
	mux := http.NewServeMux()
	cfg := apiConfig{}

	handler := http.StripPrefix("/app/", http.FileServer(http.Dir(".")))

	mux.Handle("/app/", cfg.middlewareMetricsInc(handler))
	mux.Handle("/app/assets/", http.StripPrefix("/app/assets/", http.FileServer(http.Dir("./assets"))))
	mux.HandleFunc("/healthz", handlerReadiness)
	mux.HandleFunc("/metrics", cfg.metricHandler)
	mux.HandleFunc("/reset", cfg.resetHandler)

	server := &http.Server{
		Handler: mux,
		Addr:    ":" + port,
	}

	err := server.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}

func handlerReadiness(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(http.StatusText(http.StatusOK)))
}
