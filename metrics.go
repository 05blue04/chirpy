package main

import (
	"fmt"
	"net/http"
)

func (cfg *apiConfig) metricHandler(w http.ResponseWriter, r *http.Request) {
	s := fmt.Sprintf(`
<html>
  <body>
    <h1>Welcome, Chirpy Admin</h1>
    <p>Chirpy has been visited %d times!</p>
  </body>
</html>`, cfg.fileserverHits.Load())

	w.Header().Add("Content-Type", "text/html")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(s))
}
