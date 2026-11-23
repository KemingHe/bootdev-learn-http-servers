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

func (cfg *apiConfig) middelwareMitricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.fileserverHits.Add(1)
		next.ServeHTTP(w, r)
	})
}

func handleReadiness(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(http.StatusText(http.StatusOK)))
}

func (cfg *apiConfig) handleMetricsLookup(
	w http.ResponseWriter,
	req *http.Request,
) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	data := fmt.Sprintf("Hits: %d", cfg.fileserverHits.Load())
	w.Write([]byte(data))
}

func (cfg *apiConfig) handleMetricsReset(
	w http.ResponseWriter,
	req *http.Request,
) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	cfg.fileserverHits.Store(0)
	data := fmt.Sprintf(
		"Reset Successful!\nHits: %d",
		cfg.fileserverHits.Load(),
	)
	w.Write([]byte(data))
}

func handleFileServe(fsPref, fsRoot string) http.Handler {
	return http.StripPrefix("/app", http.FileServer(http.Dir(fsRoot)))
}

func main() {
	const fsPref = "/app"
	const fsRoot = "."
	const port = "8080"

	var count atomic.Int32 // Init to zero val of 0
	apiCfg := apiConfig{fileserverHits: count}

	sMux := http.NewServeMux()
	sMux.HandleFunc("/healthz", handleReadiness)
	sMux.HandleFunc("/metrics", apiCfg.handleMetricsLookup)
	sMux.HandleFunc("/reset", apiCfg.handleMetricsReset)
	sMux.Handle(
		fsPref+"/",
		apiCfg.middelwareMitricsInc(handleFileServe(fsPref, fsRoot)),
	)

	serv := &http.Server{
		Addr:    ":" + port,
		Handler: sMux,
	}

	log.Printf("Serving files from %s on port %s\n", fsRoot, port)
	log.Fatal(serv.ListenAndServe()) // Start server
}
