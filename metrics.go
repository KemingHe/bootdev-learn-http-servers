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

func (cfg *apiConfig) middlewareMitricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.fileserverHits.Add(1)
		log.Printf(
			"Hit #%d: %s %s",
			cfg.fileserverHits.Load(),
			r.Method,
			r.URL.Path,
		)
		next.ServeHTTP(w, r)
	})
}

func (cfg *apiConfig) handleMetricsLookup(
	w http.ResponseWriter,
	req *http.Request,
) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	data := fmt.Sprintf(`<html>
	<body>
    <h1>Welcome, Chirpy Admin</h1>
    <p>Chirpy has been visited %d times!</p>
  </body>
</html>`, cfg.fileserverHits.Load())
	w.Write([]byte(data))
}

func (cfg *apiConfig) handleMetricsReset(
	w http.ResponseWriter,
	req *http.Request,
) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	cfg.fileserverHits.Store(0)
	log.Printf("Reset fileserverHits to %d", cfg.fileserverHits.Load())
	data := fmt.Sprintf(
		"Reset Successful!\nHits: %d",
		cfg.fileserverHits.Load(),
	)
	w.Write([]byte(data))
}
